package articlelogic

import (
	"rpc-content/pkg/code"
	"context"
	"sort"
	"strconv"
	"sync"
	"time"

	"rpc-content/client/user/user"
	"rpc-content/content"
	model "rpc-content/internal/model/article"
	"rpc-content/internal/svc"

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

		articles, err := l.svcCtx.ArticleModel.SearchArticles(l.ctx, in.Keyword, authorIds, 2, int(in.PageSize)+1, in.Cursor)
		if err != nil {
			resp.Code = int64(code.ServerErr.Code())
			resp.Msg = err.Error()
			return resp, nil
		}

		if len(articles) > int(in.PageSize) {
			articles = articles[:in.PageSize]
		} else {
			isEnd = true
		}

		items = articlesToItems(articles)

		if len(items) > 0 && !isEnd {
			last := items[len(items)-1]
			t, err := time.ParseInLocation("2006-01-02 15:04:05", last.PublishTime, time.Local)
			if err == nil {
				cursor = t.Unix()
			}
		}
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
