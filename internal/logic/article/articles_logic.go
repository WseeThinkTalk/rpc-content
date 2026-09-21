package articlelogic

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"rpc-content/content"
	model "rpc-content/internal/model/article"
	"rpc-content/internal/svc"
	types "rpc-content/internal/types/article"
	"rpc-content/pkg/code"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

type ArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticlesLogic {
	return &ArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticlesLogic) Articles(in *content.ArticlesRequest) (resp *content.ArticlesResponse, err error) {
	resp = new(content.ArticlesResponse)
	resp.Data = new(content.ArticlesData)
	resp.Data.Articles = make([]*content.ArticleItem, 0)

	if in.SortType != types.SortPublishTime && in.SortType != types.SortLikeCount {
		resp.Code = int64(code.SortTypeInvalid.Code())
		resp.Msg = code.SortTypeInvalid.Message()
		return resp, nil
	}
	if in.UserId <= 0 {
		resp.Code = int64(code.UserIdInvalid.Code())
		resp.Msg = code.UserIdInvalid.Message()
		return resp, nil
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}

	if in.Cursor == 0 {
		if in.SortType == types.SortLikeCount {
			in.Cursor = 9999999999
		} else {
			in.Cursor = time.Now().Unix()
		}
	}

	var (
		isCache, isEnd bool
		lastId, cursor int64
		articleIds     []int64
		articleModels  []*model.Article
		curPageIds     []int64
	)

	key := articlesKey(in.UserId, in.SortType)
	b, _ := l.svcCtx.BizRedis.ExistsCtx(l.ctx, key)

	if b {
		isCache = true
		_ = l.svcCtx.BizRedis.ExpireCtx(l.ctx, key, 3600*24*2)
		pairs, err := l.svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(l.ctx, key, 0, in.Cursor, 0, int(in.PageSize))
		if err != nil {
			resp.Code = int64(code.ServerErr.Code())
			resp.Msg = err.Error()
			return resp, nil
		}

		// 解析 Redis ZSet 成员及对应分值
		var scores []int64
		for _, v := range pairs {
			artId, err := strconv.ParseInt(v.Key, 10, 64)
			if err != nil {
				continue
			}
			score := v.Score
			articleIds = append(articleIds, artId)
			scores = append(scores, score)
		}

		if len(articleIds) > 0 {
			if articleIds[len(articleIds)-1] == -1 {
				articleIds = articleIds[:len(articleIds)-1]
				scores = scores[:len(scores)-1]
				isEnd = true
			}
		}

		curPageIds = articleIds

		if in.Cursor > 0 && in.ArticleId > 0 && len(curPageIds) > 0 {
			// 定位游标在当前缓存ID列表中的偏移量
			for i, v := range curPageIds {
				if scores[i] == in.Cursor && v == in.ArticleId {
					curPageIds = curPageIds[i+1:]
					break
				}
			}
		}
	} else {
		articleModels, err = l.svcCtx.ArticleModel.ArticlesByUserIdWithoutCursor(l.ctx, in.UserId, -1)
		if err != nil {
			resp.Code = int64(code.ServerErr.Code())
			resp.Msg = err.Error()
			return resp, nil
		}

		if len(articleModels) == 0 {
			resp.Data.IsEnd = true
			return resp, nil
		}

		// 根据游标过滤早于指定分值的文章
		var filtered []*model.Article
		for _, v := range articleModels {
			var score int64
			if in.SortType == types.SortLikeCount {
				score = v.LikeNum
			} else {
				score = v.PublishTime.Unix()
			}
			if score <= in.Cursor {
				filtered = append(filtered, v)
			}
		}

		if in.Cursor > 0 && in.ArticleId > 0 && len(filtered) > 0 {
			// 定位数据库结果中游标文章的具体位置
			for i, v := range filtered {
				var score int64
				if in.SortType == types.SortLikeCount {
					score = v.LikeNum
				} else {
					score = v.PublishTime.Unix()
				}
				if score == in.Cursor && v.Id == in.ArticleId {
					filtered = filtered[i+1:]
					break
				}
			}
		}

		var firstPage []*model.Article
		if len(filtered) > int(in.PageSize) {
			firstPage = filtered[:in.PageSize]
			isEnd = false
		} else {
			firstPage = filtered
			isEnd = true
		}

		// 收集分页内的文章ID
		for _, v := range firstPage {
			curPageIds = append(curPageIds, v.Id)
		}
	}

	// 批量查询文章详情并构建响应列表
	items := make([]*content.ArticleItem, 0, len(curPageIds))
	for _, v := range curPageIds {
		art, err := l.svcCtx.ArticleModel.FindOne(l.ctx, v)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				continue
			}
			resp.Code = int64(code.ServerErr.Code())
			resp.Msg = err.Error()
			return resp, nil
		}
		items = append(items, &content.ArticleItem{
			Id:           art.Id,
			Title:        art.Title,
			Content:      art.Content,
			Description:  art.Description,
			Cover:        art.Cover,
			AuthorId:     art.AuthorId,
			LikeCount:    art.LikeNum,
			CommentCount: art.CommentNum,
			PublishTime:  art.PublishTime.Unix(),
			Status:       int64(art.Status),
		})
	}

	if len(items) > 0 {
		last := items[len(items)-1]
		lastId = last.Id
		if in.SortType == types.SortLikeCount {
			cursor = last.LikeCount
		} else {
			cursor = last.PublishTime
		}
		if cursor < 0 {
			cursor = 0
		}
	}

	resp.Data.Articles = items
	resp.Data.IsEnd = isEnd
	resp.Data.Cursor = cursor
	resp.Data.ArticleId = lastId

	if !isCache {
		threading.GoSafe(func() {
			_ = l.addCacheArticles(context.Background(), in.UserId, articleModels)
		})
	}

	return resp, nil
}

func (l *ArticlesLogic) addCacheArticles(ctx context.Context, userId int64, articles []*model.Article) error {
	publishTimeKey := articlesKey(userId, types.SortPublishTime)
	likeNumKey := articlesKey(userId, types.SortLikeCount)

	if len(articles) == 0 {
		_, _ = l.svcCtx.BizRedis.ZaddCtx(ctx, publishTimeKey, 0, "-1")
		_, _ = l.svcCtx.BizRedis.ZaddCtx(ctx, likeNumKey, 0, "-1")
	} else {
		// 批量将文章写入 Redis 排序集合
		for _, v := range articles {
			artIdStr := strconv.FormatInt(v.Id, 10)
			_, _ = l.svcCtx.BizRedis.ZaddCtx(ctx, publishTimeKey, v.PublishTime.Unix(), artIdStr)
			_, _ = l.svcCtx.BizRedis.ZaddCtx(ctx, likeNumKey, v.LikeNum, artIdStr)
		}
		_, _ = l.svcCtx.BizRedis.ZaddCtx(ctx, publishTimeKey, 0, "-1")
		_, _ = l.svcCtx.BizRedis.ZaddCtx(ctx, likeNumKey, 0, "-1")
	}

	_ = l.svcCtx.BizRedis.ExpireCtx(ctx, publishTimeKey, 3600*24*2)
	_ = l.svcCtx.BizRedis.ExpireCtx(ctx, likeNumKey, 3600*24*2)
	return nil
}

func articlesKey(uid int64, sortType int32) string {
	return fmt.Sprintf("biz#articles#%d#%d", uid, sortType)
}
