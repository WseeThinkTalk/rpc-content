package taglogic

import (
	"context"

	"rpc-content/content"
	"rpc-content/internal/svc"
	"rpc-content/pkg/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type UntagResourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUntagResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UntagResourceLogic {
	return &UntagResourceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UntagResourceLogic) UntagResource(in *content.UntagResourceRequest) (resp *content.UntagResourceResponse, err error) {
	resp = new(content.UntagResourceResponse)

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

	exist, err := l.svcCtx.TagResourceModel.FindByTagIDAndBizIDAndTargetID(l.ctx, in.TagId, in.BizId, in.TargetId)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}
	if exist == nil {
		return resp, nil
	}

	err = l.svcCtx.TagResourceModel.Delete(l.ctx, exist.ID)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}

	return resp, nil
}
