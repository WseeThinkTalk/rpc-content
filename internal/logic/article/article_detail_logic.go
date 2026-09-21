package articlelogic

import (
	"context"
	"errors"

	"rpc-content/content"
	"rpc-content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ArticleDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleDetailLogic {
	return &ArticleDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleDetailLogic) ArticleDetail(in *content.ArticleDetailRequest) (resp *content.ArticleDetailResponse, err error) {
	resp = new(content.ArticleDetailResponse)
	resp.Code = 200
	resp.Msg = "success"
	resp.Data = new(content.ArticleItem)

	article, err := l.svcCtx.ArticleModel.FindOne(l.ctx, in.ArticleId)
	if err != nil {
		resp.Data = nil
		if errors.Is(err, sqlx.ErrNotFound) {
			resp.Code = 404
			resp.Msg = "文章不存在"
			return resp, nil
		}
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}
	if article == nil {
		resp.Data = nil
		resp.Code = 404
		resp.Msg = "文章不存在"
		return resp, nil
	}

	resp.Data.Id = article.Id
	resp.Data.Title = article.Title
	resp.Data.Content = article.Content
	resp.Data.Description = article.Description
	resp.Data.Cover = article.Cover
	resp.Data.AuthorId = article.AuthorId
	resp.Data.LikeCount = article.LikeNum
	resp.Data.CommentCount = article.CommentNum
	resp.Data.PublishTime = article.PublishTime.Unix()
	resp.Data.Status = int64(article.Status)

	return resp, nil
}
