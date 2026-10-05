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
	"rpc-content/pkg/sensitive"
	"rpc-content/pkg/snowflake"

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

	// 敏感词检查
	defaultFilter := sensitive.NewFilter([]string{"涉黄", "涉暴", "赌博", "违禁"})
	if defaultFilter.IsSensitive(in.Title) || defaultFilter.IsSensitive(in.Description) || defaultFilter.IsSensitive(in.Content) {
		resp.Code = int64(code.ContentContainsSensitiveWord.Code())
		resp.Msg = code.ContentContainsSensitiveWord.Message()
		return resp, nil
	}

	// 预分配分布式唯一文章ID
	articleId := snowflake.GenerateID()

	_, err = l.svcCtx.ArticleModel.Insert(l.ctx, &model.Article{
		Id:          articleId,
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
		resp.Code = int64(code.ServerErr.Code())
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

	// 关联文章标签
	if len(in.TagIds) > 0 {
		for _, v := range in.TagIds {
			_ = l.svcCtx.TagResourceModel.Insert(l.ctx, &tagmodel.TagResource{
				BizID:      "article",
				TargetID:   articleId,
				TagID:      v,
				UserID:     in.UserId,
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			})
		}
	}

	resp.Data.ArticleId = articleId
	return resp, nil
}
