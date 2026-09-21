package articlelogic

import (
	"context"

	"rpc-content/content"
	"rpc-content/internal/svc"
	types "rpc-content/internal/types/article"
	"rpc-content/pkg/code"
	"rpc-content/pkg/xcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleDeleteLogic {
	return &ArticleDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleDeleteLogic) ArticleDelete(in *content.ArticleDeleteRequest) (resp *content.ArticleDeleteResponse, err error) {
	resp = new(content.ArticleDeleteResponse)

	if in.UserId <= 0 {
		resp.Code = int64(code.UserIdInvalid.Code())
		resp.Msg = code.UserIdInvalid.Message()
		return resp, nil
	}
	if in.ArticleId <= 0 {
		resp.Code = int64(code.ArticleIdInvalid.Code())
		resp.Msg = code.ArticleIdInvalid.Message()
		return resp, nil
	}
	article, err := l.svcCtx.ArticleModel.FindOne(l.ctx, in.ArticleId)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}
	if article == nil {
		resp.Code = 404
		resp.Msg = "文章不存在"
		return resp, nil
	}
	if article.AuthorId != in.UserId {
		resp.Code = int64(xcode.AccessDenied.Code())
		resp.Msg = xcode.AccessDenied.Message()
		return resp, nil
	}
	err = l.svcCtx.ArticleModel.UpdateArticleStatus(l.ctx, in.ArticleId, types.ArticleStatusUserDelete)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}
	_, _ = l.svcCtx.BizRedis.ZremCtx(l.ctx, articlesKey(in.UserId, types.SortPublishTime), in.ArticleId)
	_, _ = l.svcCtx.BizRedis.ZremCtx(l.ctx, articlesKey(in.UserId, types.SortLikeCount), in.ArticleId)
	_, _ = l.svcCtx.BizRedis.ZremCtx(l.ctx, "biz#articles#global", in.ArticleId)

	return resp, nil
}
