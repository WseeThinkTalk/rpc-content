package articlelogic

import (
	"rpc-content/pkg/code"
	"context"

	"rpc-content/content"
	"rpc-content/internal/svc"
	types "rpc-content/internal/types/article"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminAuditLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminAuditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminAuditLogic {
	return &AdminAuditLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminAuditLogic) AdminAudit(in *content.AdminAuditRequest) (resp *content.AdminAuditResponse, err error) {
	resp = new(content.AdminAuditResponse)

	if int(in.Status) != types.ArticleStatusNotPass && int(in.Status) != types.ArticleStatusVisible {
		resp.Code = int64(code.RequestErr.Code())
		resp.Msg = "无效的审核状态，仅允许 1=拒绝 或 2=通过"
		return resp, nil
	}

	article, err := l.svcCtx.ArticleModel.FindOne(l.ctx, in.ArticleId)
	if err != nil {
		resp.Code = int64(code.ServerErr.Code())
		resp.Msg = err.Error()
		return resp, nil
	}
	if article == nil {
		resp.Code = int64(code.NotFound.Code())
		resp.Msg = "文章不存在"
		return resp, nil
	}
	if article.Status != types.ArticleStatusPending {
		resp.Code = int64(code.RequestErr.Code())
		resp.Msg = "仅待审核状态的文章可以审核"
		return resp, nil
	}

	err = l.svcCtx.ArticleModel.UpdateArticleStatus(l.ctx, in.ArticleId, int(in.Status))
	if err != nil {
		resp.Code = int64(code.ServerErr.Code())
		resp.Msg = err.Error()
		return resp, nil
	}

	if int(in.Status) == types.ArticleStatusVisible {
		_, _ = l.svcCtx.BizRedis.DelCtx(l.ctx, "biz#articles#global")
	}

	return resp, nil
}
