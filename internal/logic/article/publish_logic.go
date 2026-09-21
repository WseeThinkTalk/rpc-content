package articlelogic

import (
	"context"
	"strconv"
	"time"

	"rpc-content/content"
	model "rpc-content/internal/model/article"
	tagmodel "rpc-content/internal/model/tag"
	"rpc-content/internal/svc"
	types "rpc-content/internal/types/article"
	"rpc-content/pkg/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

const (
	maxTitleLength       = 200
	maxContentLength     = 100000
	maxDescriptionLength = 500
)

func (l *PublishLogic) Publish(in *content.PublishRequest) (resp *content.PublishResponse, err error) {
	resp = new(content.PublishResponse)
	resp.Data = new(content.PublishData)

	if in.UserId <= 0 {
		resp.Code = int64(code.UserIdInvalid.Code())
		resp.Msg = code.UserIdInvalid.Message()
		return resp, nil
	}
	if len(in.Title) == 0 {
		resp.Code = int64(code.ArticleTitleCantEmpty.Code())
		resp.Msg = code.ArticleTitleCantEmpty.Message()
		return resp, nil
	}
	if len(in.Title) > maxTitleLength {
		resp.Code = int64(code.ArticleTitleTooLong.Code())
		resp.Msg = code.ArticleTitleTooLong.Message()
		return resp, nil
	}
	if len(in.Content) == 0 {
		resp.Code = int64(code.ArticleContentCantEmpty.Code())
		resp.Msg = code.ArticleContentCantEmpty.Message()
		return resp, nil
	}
	if len(in.Content) > maxContentLength {
		resp.Code = int64(code.ArticleContentTooLong.Code())
		resp.Msg = code.ArticleContentTooLong.Message()
		return resp, nil
	}
	if len(in.Description) > maxDescriptionLength {
		resp.Code = int64(code.ArticleDescTooLong.Code())
		resp.Msg = code.ArticleDescTooLong.Message()
		return resp, nil
	}

	ret, err := l.svcCtx.ArticleModel.Insert(l.ctx, &model.Article{
		AuthorId:    in.UserId,
		Title:       in.Title,
		Content:     in.Content,
		Description: in.Description,
		Cover:       in.Cover,
		Status:      types.ArticleStatusPending, // 待审核
		PublishTime: time.Now(),
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	})
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}

	articleId, err := ret.LastInsertId()
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}

	var (
		articleIdStr   = strconv.FormatInt(articleId, 10)
		publishTimeKey = articlesKey(in.UserId, types.SortPublishTime)
		likeNumKey     = articlesKey(in.UserId, types.SortLikeCount)
	)
	b, _ := l.svcCtx.BizRedis.ExistsCtx(l.ctx, publishTimeKey)
	if b {
		_, _ = l.svcCtx.BizRedis.ZaddCtx(l.ctx, publishTimeKey, time.Now().Unix(), articleIdStr)
	}
	b, _ = l.svcCtx.BizRedis.ExistsCtx(l.ctx, likeNumKey)
	if b {
		_, _ = l.svcCtx.BizRedis.ZaddCtx(l.ctx, likeNumKey, 0, articleIdStr)
	}

	// 关联标签
	if len(in.TagIds) > 0 {
		for _, tagId := range in.TagIds {
			_ = l.svcCtx.TagResourceModel.Insert(l.ctx, &tagmodel.TagResource{
				BizID:      "article",
				TargetID:   articleId,
				TagID:      tagId,
				UserID:     in.UserId,
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			})
		}
	}

	resp.Data.ArticleId = articleId
	return resp, nil
}
