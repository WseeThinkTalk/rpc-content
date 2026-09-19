package logic

import (
	"context"
	"sort"
	"strconv"
	"sync"
	"time"

	"rpc-content/article/internal/model"
	"rpc-content/article/internal/svc"
	"rpc-content/article/pb"
	"rpc-content/client/user/user"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
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

func (l *SearchArticlesLogic) SearchArticles(in *pb.SearchRequest) (*pb.SearchResponse, error) {
	if in.PageSize == 0 {
		in.PageSize = 20
	}

	var (
		items  []*pb.SearchItem
		isEnd  bool
		cursor int64
	)

	if in.Keyword == "" {
		// 全量文章 — 优先从 Redis 拿 ID，否则从 DB 全量加载
		var allModels []*model.Article

		if in.SortType == 1 {
			// 按点赞排序
			var err error
			allModels, err = l.svcCtx.ArticleMOdel.ArticlesAllVisible(l.ctx, 2)
			if err != nil {
				l.Logger.Errorf("[SearchArticles] ArticlesAllVisible error: %v", err)
				return nil, err
			}
			sort.SliceStable(allModels, func(i, j int) bool {
				if allModels[i].LikeNum == allModels[j].LikeNum {
					return allModels[i].PublishTime.After(allModels[j].PublishTime)
				}
				return allModels[i].LikeNum > allModels[j].LikeNum
			})
		} else {
			// 时间排序 — 尝试 Redis 缓存
			key := "biz#articles#global"
			cachedIds := getCachedArticleIds(l.ctx, l.svcCtx, key)

			if len(cachedIds) > 0 {
				// 缓存命中：批量查询
				models, err := l.svcCtx.ArticleMOdel.FindByIds(l.ctx, cachedIds)
				if err != nil {
					l.Logger.Errorf("[SearchArticles] FindByIds error: %v", err)
					return nil, err
				}
				// 按 publish_time desc 重排
				sort.SliceStable(models, func(i, j int) bool {
					return models[i].PublishTime.After(models[j].PublishTime)
				})
				allModels = models
			} else {
				// 缓存未命中：全量 DB 加载
				models, err := l.svcCtx.ArticleMOdel.ArticlesAllVisible(l.ctx, 2)
				if err != nil {
					l.Logger.Errorf("[SearchArticles] ArticlesAllVisible error: %v", err)
					return nil, err
				}
				allModels = models

				// 异步构建缓存
				threading.GoSafe(func() {
					_ = l.addCacheAllArticles(context.Background(), models)
				})
			}
		}

		if len(allModels) == 0 {
			return &pb.SearchResponse{IsEnd: true}, nil
		}

		// 内存分页（保证不发生重复）
		if in.Cursor == 0 {
			in.Cursor = time.Now().Unix()
		}
		var page []*model.Article
		page, isEnd, cursor = paginateByTime(allModels, in.Cursor, int(in.PageSize))
		items = articlesToItems(page)
	} else {
		// 1. 查找匹配该关键词的所有作者 ID
		var authorIds []int64
		if in.Keyword != "" {
			userResp, err := l.svcCtx.UserRPC.AdminUserList(l.ctx, &user.AdminUserListRequest{
				Keyword:  in.Keyword,
				PageSize: 100,
			})
			if err == nil && userResp != nil {
				for _, item := range userResp.Items {
					authorIds = append(authorIds, item.UserId)
				}
			}
		}

		// 2. 关键词与作者 ID 双重条件搜索
		articles, err := l.svcCtx.ArticleMOdel.SearchArticles(l.ctx, in.Keyword, authorIds, 2, int(in.PageSize)+1, in.Cursor)
		if err != nil {
			l.Logger.Errorf("[SearchArticles] SearchArticles error: %v", err)
			return nil, err
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

	// 并发批量查询作者名
	populateAuthorNames(l.ctx, l.svcCtx, items, l.Logger)

	return &pb.SearchResponse{
		Items:  items,
		Cursor: cursor,
		IsEnd:  isEnd,
	}, nil
}

// getCachedArticleIds 从 Redis 获取全部缓存文章 ID（按 publish_time desc 排序）
func getCachedArticleIds(ctx context.Context, svcCtx *svc.ServiceContext, key string) []int64 {
	exists, err := svcCtx.BizRedis.ExistsCtx(ctx, key)
	if err != nil || !exists {
		return nil
	}
	_ = svcCtx.BizRedis.ExpireCtx(ctx, key, 3600*24*2)

	// 获取全部（最多 10000 条）
	pairs, err := svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(
		ctx, key, 0, time.Now().Unix()+86400, 0, 10000)
	if err != nil {
		return nil
	}

	ids := make([]int64, 0, len(pairs))
	for _, pair := range pairs {
		id, err := strconv.ParseInt(pair.Key, 10, 64)
		if err != nil || id == -1 {
			continue
		}
		ids = append(ids, id)
	}
	return ids
}

// articlesToItems 模型 → PB 转换
func articlesToItems(models []*model.Article) []*pb.SearchItem {
	items := make([]*pb.SearchItem, 0, len(models))
	for _, art := range models {
		items = append(items, &pb.SearchItem{
			ArticleId:   art.Id,
			Title:       art.Title,
			Description: art.Description,
			Cover:       art.Cover,
			AuthorId:    art.AuthorId,
			LikeNum:     art.LikeNum,
			CommentNum:  art.CommentNum,
			PublishTime: art.PublishTime.Format("2006-01-02 15:04:05"),
		})
	}
	return items
}

// paginateByTime 全量数据按时间分页
func paginateByTime(allArticles []*model.Article, cursor int64, pageSize int) ([]*model.Article, bool, int64) {
	cursorTime := time.Unix(cursor, 0)

	// 找到严格在 cursor 之后的第一条（publish_time < cursor）
	// 用 Before 不用 Equal，避免上一页最后一条重复
	startIdx := len(allArticles)
	for i, art := range allArticles {
		if art.PublishTime.Before(cursorTime) {
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

// populateAuthorNames 并发批量查询作者名
func populateAuthorNames(ctx context.Context, svcCtx *svc.ServiceContext, items []*pb.SearchItem, log logx.Logger) {
	if len(items) == 0 {
		return
	}

	authorIds := make(map[int64]struct{}, len(items))
	for _, item := range items {
		authorIds[item.AuthorId] = struct{}{}
	}

	type result struct {
		authorId int64
		name     string
		avatar   string
	}
	results := make(chan result, len(authorIds))
	var wg sync.WaitGroup

	for uid := range authorIds {
		wg.Add(1)
		go func(userId int64) {
			defer wg.Done()
			u, err := svcCtx.UserRPC.FindById(ctx, &user.FindByIdRequest{UserId: userId})
			if err != nil {
				log.Errorf("[SearchArticles] FindById userId: %d error: %v", userId, err)
				return
			}
			results <- result{authorId: userId, name: u.Username, avatar: u.Avatar}
		}(uid)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	type authorInfo struct {
		name   string
		avatar string
	}
	infoMap := make(map[int64]authorInfo, len(authorIds))
	for r := range results {
		infoMap[r.authorId] = authorInfo{name: r.name, avatar: r.avatar}
	}

	for _, item := range items {
		if info, ok := infoMap[item.AuthorId]; ok {
			item.AuthorName = info.name
			item.AuthorAvatar = info.avatar
		}
	}
}

func (l *SearchArticlesLogic) addCacheAllArticles(ctx context.Context, articles []*model.Article) error {
	key := "biz#articles#global"
	if len(articles) == 0 {
		_, _ = l.svcCtx.BizRedis.ZaddCtx(ctx, key, 0, "-1")
	} else {
		for _, article := range articles {
			artIdStr := strconv.FormatInt(article.Id, 10)
			_, _ = l.svcCtx.BizRedis.ZaddCtx(ctx, key, article.PublishTime.Unix(), artIdStr)
		}
		_, _ = l.svcCtx.BizRedis.ZaddCtx(ctx, key, 0, "-1")
	}
	_ = l.svcCtx.BizRedis.ExpireCtx(ctx, key, 3600*24*2)
	return nil
}
