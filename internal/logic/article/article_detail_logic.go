package articlelogic

import (
	"rpc-content/pkg/code"
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

// ArticleDetail 获取文章详情
func (l *ArticleDetailLogic) ArticleDetail(in *content.ArticleDetailRequest) (resp *content.ArticleDetailResponse, err error) {
	resp = new(content.ArticleDetailResponse)
	resp.Data = new(content.ArticleItem)

	// 查询文章信息
	article, err := l.svcCtx.ArticleModel.FindOne(l.ctx, in.ArticleId)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			resp.Code = int64(code.NotFound.Code())
			resp.Msg = "文章不存在"
			return resp, nil
		}
		resp.Code = int64(code.ServerErr.Code())
		resp.Msg = err.Error()
		return resp, nil
	}
	if article == nil {
		resp.Code = int64(code.NotFound.Code())
		resp.Msg = "文章不存在"
		return resp, nil
	}

	// 组装返回结果
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
