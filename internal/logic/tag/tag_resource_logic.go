package taglogic

import (
	"context"
	"time"

	"rpc-content/content"
	model "rpc-content/internal/model/tag"
	"rpc-content/internal/svc"
	"rpc-content/pkg/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagResourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTagResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagResourceLogic {
	return &TagResourceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TagResourceLogic) TagResource(in *content.TagResourceRequest) (resp *content.TagResourceResponse, err error) {
	resp = new(content.TagResourceResponse)

	if in.BizId == "" {
		return nil, code.BizIdEmpty
	}
	if in.TargetId == 0 {
		return nil, code.TargetIdEmpty
	}
	if in.TagId == 0 {
		return nil, code.TagIdEmpty
	}
	if in.UserId == 0 {
		return nil, code.TagUserIdEmpty
	}

	tag, err := l.svcCtx.TagModel.FindOne(l.ctx, in.TagId)
	if err != nil {
		l.Logger.Errorf("[TagResource] TagModel.FindOne err: %v tagId: %d", err, in.TagId)
		return nil, err
	}
	if tag == nil {
		return nil, code.TagNotFound
	}

	exist, err := l.svcCtx.TagResourceModel.FindByTagIDAndBizIDAndTargetID(l.ctx, in.TagId, in.BizId, in.TargetId)
	if err != nil {
		l.Logger.Errorf("[TagResource] TagResourceModel.FindByTagIDAndBizIDAndTargetID err: %v req: %+v", err, in)
		return nil, err
	}
	if exist != nil {
		return nil, code.TagResourceExists
	}

	tr := &model.TagResource{
		BizID:      in.BizId,
		TargetID:   in.TargetId,
		TagID:      in.TagId,
		UserID:     in.UserId,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	if err := l.svcCtx.TagResourceModel.Insert(l.ctx, tr); err != nil {
		l.Logger.Errorf("[TagResource] TagResourceModel.Insert err: %v tr: %+v", err, tr)
		return nil, err
	}

	return resp, nil
}
