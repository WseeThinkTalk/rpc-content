package logic

import (
	"context"

	"rpc-content/tag/code"
	"rpc-content/tag/internal/svc"
	"rpc-content/tag/pb"

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

func (l *UntagResourceLogic) UntagResource(in *pb.UntagResourceRequest) (*pb.UntagResourceResponse, error) {
	if in.BizId == "" {
		return nil, code.BizIdEmpty
	}
	if in.TargetId == 0 {
		return nil, code.TargetIdEmpty
	}
	if in.TagId == 0 {
		return nil, code.TagIdEmpty
	}

	exist, err := l.svcCtx.TagResourceModel.FindByTagIDAndBizIDAndTargetID(l.ctx, in.TagId, in.BizId, in.TargetId)
	if err != nil {
		l.Logger.Errorf("[UntagResource] TagResourceModel.FindByTagIDAndBizIDAndTargetID err: %v req: %+v", err, in)
		return nil, err
	}
	if exist == nil {
		return &pb.UntagResourceResponse{}, nil
	}

	err = l.svcCtx.TagResourceModel.Delete(l.ctx, exist.ID)
	if err != nil {
		l.Logger.Errorf("[UntagResource] TagResourceModel.Delete err: %v id: %d", err, exist.ID)
		return nil, err
	}

	return &pb.UntagResourceResponse{}, nil
}
