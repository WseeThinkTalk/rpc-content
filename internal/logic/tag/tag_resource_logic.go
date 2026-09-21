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
		resp.Code = int64(code.BizIdEmpty.Code())
		resp.Msg = code.BizIdEmpty.Message()
		return resp, nil
	}
	if in.TargetId == 0 {
		resp.Code = int64(code.TargetIdEmpty.Code())
		resp.Msg = code.TargetIdEmpty.Message()
		return resp, nil
	}
	if in.TagId == 0 {
		resp.Code = int64(code.TagIdEmpty.Code())
		resp.Msg = code.TagIdEmpty.Message()
		return resp, nil
	}
	if in.UserId == 0 {
		resp.Code = int64(code.TagUserIdEmpty.Code())
		resp.Msg = code.TagUserIdEmpty.Message()
		return resp, nil
	}

	tag, err := l.svcCtx.TagModel.FindOne(l.ctx, in.TagId)
	if err != nil {
		resp.Code = int64(code.ServerErr.Code())
		resp.Msg = err.Error()
		return resp, nil
	}
	if tag == nil {
		resp.Code = int64(code.TagNotFound.Code())
		resp.Msg = code.TagNotFound.Message()
		return resp, nil
	}

	exist, err := l.svcCtx.TagResourceModel.FindByTagIDAndBizIDAndTargetID(l.ctx, in.TagId, in.BizId, in.TargetId)
	if err != nil {
		resp.Code = int64(code.ServerErr.Code())
		resp.Msg = err.Error()
		return resp, nil
	}
	if exist != nil {
		resp.Code = int64(code.TagResourceExists.Code())
		resp.Msg = code.TagResourceExists.Message()
		return resp, nil
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
		resp.Code = int64(code.ServerErr.Code())
		resp.Msg = err.Error()
		return resp, nil
	}

	return resp, nil
}
