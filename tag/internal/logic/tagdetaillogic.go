package logic

import (
	"context"

	"rpc-content/tag/code"
	"rpc-content/tag/internal/svc"
	"rpc-content/tag/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTagDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagDetailLogic {
	return &TagDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TagDetailLogic) TagDetail(in *pb.TagDetailRequest) (*pb.TagDetailResponse, error) {
	if in.TagId == 0 {
		return nil, code.TagIdEmpty
	}

	tag, err := l.svcCtx.TagModel.FindOne(l.ctx, in.TagId)
	if err != nil {
		l.Logger.Errorf("[TagDetail] TagModel.FindOne err: %v tagId: %d", err, in.TagId)
		return nil, err
	}
	if tag == nil {
		return nil, code.TagNotFound
	}

	resourceCount, err := l.svcCtx.TagResourceModel.CountByTagID(l.ctx, in.TagId)
	if err != nil {
		l.Logger.Errorf("[TagDetail] TagResourceModel.CountByTagID err: %v tagId: %d", err, in.TagId)
		return nil, err
	}

	return &pb.TagDetailResponse{
		TagId:         tag.ID,
		TagName:       tag.TagName,
		TagDesc:       tag.TagDesc,
		ResourceCount: resourceCount,
		CreateTime:    tag.CreateTime.Unix(),
	}, nil
}
