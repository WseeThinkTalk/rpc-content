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

func (l *ArticleDetailLogic) ArticleDetail(in *content.ArticleDetailRequest) (*content.ArticleDetailResponse, error) {
	article, err := l.svcCtx.ArticleModel.FindOne(l.ctx, in.ArticleId)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return &content.ArticleDetailResponse{
				Code: 200,
				Msg:  "success",
			}, nil
		}
		return nil, err
	}
	return &content.ArticleDetailResponse{
		Code: 200,
		Msg:  "success",
		Data: &content.ArticleItem{
			Id:           article.Id,
			Title:        article.Title,
			Content:      article.Content,
			Description:  article.Description,
			Cover:        article.Cover,
			AuthorId:     article.AuthorId,
			LikeCount:    article.LikeNum,
			CommentCount: article.CommentNum,
			PublishTime:  article.PublishTime.Unix(),
			Status:       int64(article.Status),
		},
	}, nil
}
