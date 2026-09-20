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
	resp.Code = 200
	resp.Msg = "success"
	resp.Data = new(content.PublishData)

	if in.UserId <= 0 {
		return nil, code.UserIdInvalid
	}
	if len(in.Title) == 0 {
		return nil, code.ArticleTitleCantEmpty
	}
	if len(in.Title) > maxTitleLength {
		return nil, code.ArticleTitleTooLong
	}
	if len(in.Content) == 0 {
		return nil, code.ArticleContentCantEmpty
	}
	if len(in.Content) > maxContentLength {
		return nil, code.ArticleContentTooLong
	}
	if len(in.Description) > maxDescriptionLength {
		return nil, code.ArticleDescTooLong
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
		l.Logger.Errorf("Publish Insert req: %v error: %v", in, err)
		return nil, err
	}

	articleId, err := ret.LastInsertId()
	if err != nil {
		l.Logger.Errorf("LastInsertId error: %v", err)
		return nil, err
	}

	var (
		articleIdStr   = strconv.FormatInt(articleId, 10)
		publishTimeKey = articlesKey(in.UserId, types.SortPublishTime)
		likeNumKey     = articlesKey(in.UserId, types.SortLikeCount)
	)
	b, _ := l.svcCtx.BizRedis.ExistsCtx(l.ctx, publishTimeKey)
	if b {
		_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, publishTimeKey, time.Now().Unix(), articleIdStr)
		if err != nil {
			logx.Errorf("ZaddCtx req: %v error: %v", in, err)
		}
	}
	b, _ = l.svcCtx.BizRedis.ExistsCtx(l.ctx, likeNumKey)
	if b {
		_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, likeNumKey, 0, articleIdStr)
		if err != nil {
			logx.Errorf("ZaddCtx req: %v error: %v", in, err)
		}
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
