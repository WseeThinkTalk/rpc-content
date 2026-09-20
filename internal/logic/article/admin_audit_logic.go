package articlelogic

import (
	"context"
	"errors"

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

func (l *AdminAuditLogic) AdminAudit(in *content.AdminAuditRequest) (*content.AdminAuditResponse, error) {
	if int(in.Status) != types.ArticleStatusNotPass && int(in.Status) != types.ArticleStatusVisible {
		l.Logger.Errorf("[AdminAudit] invalid status: %d, articleId: %d", in.Status, in.ArticleId)
		return nil, errors.New("无效的审核状态，仅允许 1=拒绝 或 2=通过")
	}

	article, err := l.svcCtx.ArticleModel.FindOne(l.ctx, in.ArticleId)
	if err != nil {
		l.Logger.Errorf("[AdminAudit] FindById error: %v, articleId: %d", err, in.ArticleId)
		return nil, err
	}
	if article == nil {
		return nil, errors.New("文章不存在")
	}
	if article.Status != types.ArticleStatusPending {
		return nil, errors.New("仅待审核状态的文章可以审核")
	}

	err = l.svcCtx.ArticleModel.UpdateArticleStatus(l.ctx, in.ArticleId, int(in.Status))
	if err != nil {
		l.Logger.Errorf("[AdminAudit] UpdateArticleStatus error: %v, articleId: %d, status: %d", err, in.ArticleId, in.Status)
		return nil, err
	}

	if int(in.Status) == types.ArticleStatusVisible {
		_, err = l.svcCtx.BizRedis.DelCtx(l.ctx, "biz#articles#global")
		if err != nil {
			l.Logger.Errorf("[AdminAudit] DelCtx biz#articles#global error: %v", err)
		}
	}

	return &content.AdminAuditResponse{
		Code: 200,
		Msg:  "success",
	}, nil
}
