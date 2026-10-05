package articlelogic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"rpc-content/client/user/user"
	"rpc-content/content"
	model "rpc-content/internal/model/article"
	"rpc-content/internal/svc"
	"rpc-content/pkg/breaker"
	"rpc-content/pkg/cacheguard"
	"rpc-content/pkg/code"
	"rpc-content/pkg/esquery"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchArticlesLogic {
	return &SearchArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchArticlesLogic) SearchArticles(in *content.SearchRequest) (resp *content.SearchResponse, err error) {
	resp = new(content.SearchResponse)
	resp.Data = new(content.SearchData)
	resp.Data.Items = make([]*content.SearchItem, 0)

	if in.PageSize == 0 {
		in.PageSize = 20
	}

	var (
		items  []*content.SearchItem
		isEnd  bool
		cursor int64
	)

	if in.Keyword == "" {
		if in.SortType == 1 {
			// 按点赞数排序（取前200条热榜）
			allModels, err := l.svcCtx.ArticleModel.ArticlesByCursor(l.ctx, 2, time.Time{}, 200)
			if err != nil {
				resp.Code = int64(code.ServerErr.Code())
				resp.Msg = err.Error()
				return resp, nil
			}
			sort.SliceStable(allModels, func(i, j int) bool {
				if allModels[i].LikeNum == allModels[j].LikeNum {
					return allModels[i].PublishTime.After(allModels[j].PublishTime)
				}
				return allModels[i].LikeNum > allModels[j].LikeNum
			})
			if len(allModels) == 0 {
				resp.Data.IsEnd = true
				return resp, nil
			}
			if in.Cursor == 0 {
				in.Cursor = time.Now().Unix()
			}
			var page []*model.Article
			page, isEnd, cursor = paginateByTime(allModels, in.Cursor, int(in.PageSize))
			items = articlesToItems(page)
		} else {
			// 按发布时间游标分页
			var cursorTime time.Time
			if in.Cursor > 0 {
				cursorTime = time.Unix(in.Cursor, 0)
			}
			// 多查1条用于判断是否有下一页
			limit := int(in.PageSize) + 1
			models, err := l.svcCtx.ArticleModel.ArticlesByCursor(l.ctx, 2, cursorTime, limit)
			if err != nil {
				resp.Code = int64(code.ServerErr.Code())
				resp.Msg = err.Error()
				return resp, nil
			}

			// 截取当前页数据并更新下一页游标
			if len(models) > int(in.PageSize) {
				models = models[:in.PageSize]
				isEnd = false
				cursor = models[len(models)-1].PublishTime.Unix()
			} else {
				// 到达最后一页
				isEnd = true
				if len(models) > 0 {
					cursor = models[len(models)-1].PublishTime.Unix()
				}
			}
			items = articlesToItems(models)
		}
	} else {
		var authorIds []int64
		if in.Keyword != "" && l.svcCtx.UserRPC != nil {
			userResp, err := l.svcCtx.UserRPC.AdminUserList(l.ctx, &user.AdminUserListRequest{
				Keyword:  in.Keyword,
				PageSize: 100,
			})
			if err == nil && userResp != nil {
				// 收集匹配关键词的用户ID
				for _, v := range userResp.Items {
					authorIds = append(authorIds, v.UserId)
				}
			}
		}

		var (
			articles   []*model.Article
			esSearchOK bool
		)
		if l.svcCtx.Es != nil {
			articles, isEnd, cursor, esSearchOK = l.searchTwoPhase(in, authorIds)
		}

		// ES 未配置或发生异常时，平滑降级至数据库模糊查询
		if !esSearchOK {
			dbArticles, err := l.svcCtx.ArticleModel.SearchArticles(l.ctx, in.Keyword, authorIds, 2, int(in.PageSize)+1, in.Cursor)
			if err != nil {
				resp.Code = int64(code.ServerErr.Code())
				resp.Msg = err.Error()
				return resp, nil
			}

			if len(dbArticles) > int(in.PageSize) {
				articles = dbArticles[:in.PageSize]
				isEnd = false
			} else {
				articles = dbArticles
				isEnd = true
			}

			if len(articles) > 0 && !isEnd {
				last := articles[len(articles)-1]
				cursor = last.PublishTime.Unix()
			}
		}

		items = articlesToItems(articles)
	}

	if l.svcCtx.UserRPC != nil {
		populateAuthorNames(l.ctx, l.svcCtx, items, l.Logger)
	}

	resp.Data.Items = items
	resp.Data.Cursor = cursor
	resp.Data.IsEnd = isEnd
	return resp, nil
}

func getCachedArticleIds(ctx context.Context, svcCtx *svc.ServiceContext, key string) []int64 {
	exists, err := svcCtx.BizRedis.ExistsCtx(ctx, key)
	if err != nil || !exists {
		return nil
	}
	_ = svcCtx.BizRedis.ExpireCtx(ctx, key, 3600*24*2)

	pairs, err := svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(
		ctx, key, 0, time.Now().Unix()+86400, 0, 10000)
	if err != nil {
		return nil
	}

	// 解析缓存中的文章ID
	ids := make([]int64, 0, len(pairs))
	for _, v := range pairs {
		id, err := strconv.ParseInt(v.Key, 10, 64)
		if err != nil || id == -1 {
			continue
		}
		ids = append(ids, id)
	}
	return ids
}

func articlesToItems(models []*model.Article) []*content.SearchItem {
	// 转换文章实体为搜索数据项
	items := make([]*content.SearchItem, 0, len(models))
	for _, v := range models {
		items = append(items, &content.SearchItem{
			ArticleId:   v.Id,
			Title:       v.Title,
			Description: v.Description,
			Cover:       v.Cover,
			AuthorId:    v.AuthorId,
			LikeNum:     v.LikeNum,
			CommentNum:  v.CommentNum,
			PublishTime: v.PublishTime.Format("2006-01-02 15:04:05"),
		})
	}
	return items
}

func paginateByTime(allArticles []*model.Article, cursor int64, pageSize int) ([]*model.Article, bool, int64) {
	cursorTime := time.Unix(cursor, 0)

	// 定位首个早于游标时间的文章索引
	startIdx := len(allArticles)
	for i, v := range allArticles {
		if v.PublishTime.Before(cursorTime) {
			startIdx = i
			break
		}
	}

	if startIdx >= len(allArticles) {
		return nil, true, 0
	}

	endIdx := startIdx + pageSize
	isEnd := false
	if endIdx >= len(allArticles) {
		endIdx = len(allArticles)
		isEnd = true
	}

	page := allArticles[startIdx:endIdx]

	var nextCursor int64
	if !isEnd && len(page) > 0 {
		nextCursor = page[len(page)-1].PublishTime.Unix()
	}

	return page, isEnd, nextCursor
}

func populateAuthorNames(ctx context.Context, svcCtx *svc.ServiceContext, items []*content.SearchItem, log logx.Logger) {
	if len(items) == 0 {
		return
	}

	// 提取去重后的作者ID集合
	authorIds := make(map[int64]struct{}, len(items))
	for _, v := range items {
		authorIds[v.AuthorId] = struct{}{}
	}

	type result struct {
		authorId int64
		name     string
		avatar   string
	}
	results := make(chan result, len(authorIds))
	var wg sync.WaitGroup

	// 并发批量拉取作者信息
	for v := range authorIds {
		wg.Add(1)
		go func(userId int64) {
			defer wg.Done()
			u, err := svcCtx.UserRPC.FindById(ctx, &user.FindByIdRequest{UserId: userId})
			if err != nil {
				return
			}
			results <- result{authorId: userId, name: u.Username, avatar: u.Avatar}
		}(v)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	type authorInfo struct {
		name   string
		avatar string
	}
	// 汇总各协程返回的作者信息
	infoMap := make(map[int64]authorInfo, len(authorIds))
	for v := range results {
		infoMap[v.authorId] = authorInfo{name: v.name, avatar: v.avatar}
	}

	// 为搜索结果补充作者用户名与头像
	for _, v := range items {
		if info, ok := infoMap[v.AuthorId]; ok {
			v.AuthorName = info.name
			v.AuthorAvatar = info.avatar
		}
	}
}

func (l *SearchArticlesLogic) addCacheAllArticles(ctx context.Context, articles []*model.Article) error {
	key := "biz#articles#global"
	if len(articles) == 0 {
		_, _ = l.svcCtx.BizRedis.ZaddCtx(ctx, key, 0, "-1")
	} else {
		// 批量缓存文章ID及发布时间
		for _, v := range articles {
			artIdStr := strconv.FormatInt(v.Id, 10)
			_, _ = l.svcCtx.BizRedis.ZaddCtx(ctx, key, v.PublishTime.Unix(), artIdStr)
		}
		_, _ = l.svcCtx.BizRedis.ZaddCtx(ctx, key, 0, "-1")
	}
	_ = l.svcCtx.BizRedis.ExpireCtx(ctx, key, 3600*24*2)
	return nil
}

// searchTwoPhase 第一阶段：ES 仅检索 ID；第二阶段：Redis MGet 详情物化
func (l *SearchArticlesLogic) searchTwoPhase(in *content.SearchRequest, authorIds []int64) ([]*model.Article, bool, int64, bool) {
	from := 0
	limit := int(in.PageSize) + 1

	params := esquery.SearchParams{
		Keyword:   in.Keyword,
		AuthorIds: authorIds,
		Status:    2,
		From:      from,
		Size:      limit,
	}

	dsl, err := esquery.BuildArticleSearchIDOnlyDSL(params)
	if err != nil {
		l.Errorf("[searchTwoPhase] build dsl error: %v", err)
		return nil, false, 0, false
	}

	// 利用 SingleFlightGroup 阻断并发热搜击穿，并使用 Breaker 弹性熔断防护
	sfKey := fmt.Sprintf("sf:search:%s:%d:%d", in.Keyword, from, limit)
	val, err, _ := l.svcCtx.SingleFlightGroup.Do(sfKey, func() (interface{}, error) {
		return breaker.ExecuteWithFallback(l.ctx, "es:search_articles", func() ([]int64, error) {
			res, err := l.svcCtx.Es.Search(
				l.svcCtx.Es.Search.WithContext(l.ctx),
				l.svcCtx.Es.Search.WithIndex("thinktalk_article"),
				l.svcCtx.Es.Search.WithBody(strings.NewReader(dsl)),
			)
			if err != nil {
				return nil, err
			}
			defer res.Body.Close()
			if res.IsError() {
				return nil, fmt.Errorf("es search status error: %s", res.Status())
			}
			bodyBytes, err := io.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}
			ids, _, err := esquery.ParseSearchIDs(bodyBytes)
			return ids, err
		}, func(err error) ([]int64, error) {
			l.Errorf("[searchTwoPhase] ES breaker triggered fallback: %v", err)
			return nil, err
		})
	})
	if err != nil {
		l.Errorf("[searchTwoPhase] es execute error: %v", err)
		return nil, false, 0, false
	}

	ids, ok := val.([]int64)
	if !ok || len(ids) == 0 {
		return []*model.Article{}, true, 0, true
	}

	isEnd := true
	if len(ids) > int(in.PageSize) {
		ids = ids[:in.PageSize]
		isEnd = false
	}

	// 第二阶段：批量拉取缓存与回填
	articles := l.hydrateArticles(ids)
	var nextCursor int64
	if len(articles) > 0 {
		nextCursor = articles[len(articles)-1].PublishTime.Unix()
	}

	return articles, isEnd, nextCursor, true
}

// hydrateArticles 批量从缓存拉取文章详情（MGet单次RTT），缺失ID通过 SingleFlight 回源防击穿，并设空值防穿透
func (l *SearchArticlesLogic) hydrateArticles(ids []int64) []*model.Article {
	if len(ids) == 0 {
		return nil
	}

	resultMap := make(map[int64]*model.Article, len(ids))
	var missingIds []int64

	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = fmt.Sprintf("biz#article#detail:%d", id)
	}

	// 1. 批量单次 RTT MGet 详情缓存
	if l.svcCtx.BizRedis != nil {
		vals, err := l.svcCtx.BizRedis.MgetCtx(l.ctx, keys...)
		if err == nil && len(vals) == len(ids) {
			for i, val := range vals {
				id := ids[i]
				if val == "" {
					missingIds = append(missingIds, id)
					continue
				}

				var art model.Article
				if err := json.Unmarshal([]byte(val), &art); err != nil {
					missingIds = append(missingIds, id)
					continue
				}

				// 防穿透空对象检查：若为 -1 哨兵，说明数据库中确无此记录，跳过回源
				if cacheguard.IsNullArticle(art.Id) {
					continue
				}

				resultMap[id] = &art
			}
		} else {
			// Redis 异常或长度不匹配时全量标记缺失
			missingIds = append(missingIds, ids...)
		}
	} else {
		missingIds = append(missingIds, ids...)
	}

	// 2. 对未命中的 missingIds 执行 SingleFlight 保护回源，阻断突发热搜并发击穿 MySQL
	if len(missingIds) > 0 {
		sortedMissing := make([]int64, len(missingIds))
		copy(sortedMissing, missingIds)
		sort.Slice(sortedMissing, func(i, j int) bool { return sortedMissing[i] < sortedMissing[j] })

		var sb strings.Builder
		for i, mid := range sortedMissing {
			if i > 0 {
				sb.WriteByte(',')
			}
			sb.WriteString(strconv.FormatInt(mid, 10))
		}
		sfKey := fmt.Sprintf("sf:hydrate:%s", sb.String())

		val, err, _ := l.svcCtx.SingleFlightGroup.Do(sfKey, func() (interface{}, error) {
			dbArticles, err := l.svcCtx.ArticleModel.FindByIds(l.ctx, sortedMissing)
			if err != nil {
				return nil, err
			}

			foundMap := make(map[int64]*model.Article, len(dbArticles))
			for _, art := range dbArticles {
				foundMap[art.Id] = art
				if l.svcCtx.BizRedis != nil {
					data, _ := json.Marshal(art)
					_ = l.svcCtx.BizRedis.SetexCtx(l.ctx, fmt.Sprintf("biz#article#detail:%d", art.Id), string(data), cacheguard.NormalArticleTTL)
				}
			}

			// 对在 MySQL 中不存在的 ID 回填 Null Object 空值防穿透（60s 短 TTL）
			for _, mid := range sortedMissing {
				if _, ok := foundMap[mid]; !ok {
					if l.svcCtx.BizRedis != nil {
						_ = l.svcCtx.BizRedis.SetexCtx(l.ctx, fmt.Sprintf("biz#article#detail:%d", mid), cacheguard.BuildNullArticlePayload(), cacheguard.NullArticleTTL)
					}
				}
			}
			return foundMap, nil
		})

		if err == nil && val != nil {
			if foundMap, ok := val.(map[int64]*model.Article); ok {
				for id, art := range foundMap {
					resultMap[id] = art
				}
			}
		}
	}

	// 3. 严格保持原召回顺序
	ordered := make([]*model.Article, 0, len(ids))
	for _, id := range ids {
		if art, ok := resultMap[id]; ok {
			ordered = append(ordered, art)
		}
	}
	return ordered
}
