package logic

import (
	"context"

	"rpc-content/tag/code"
	"rpc-content/tag/internal/svc"
	"rpc-content/tag/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagsByResourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTagsByResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagsByResourceLogic {
	return &TagsByResourceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TagsByResourceLogic) TagsByResource(in *pb.TagsByResourceRequest) (*pb.TagsByResourceResponse, error) {
	if in.BizId == "" {
		return nil, code.BizIdEmpty
	}
	if in.TargetId == 0 {
		return nil, code.TargetIdEmpty
	}

	trs, err := l.svcCtx.TagResourceModel.FindTagsByBizIDAndTargetID(l.ctx, in.BizId, in.TargetId)
	if err != nil {
		l.Logger.Errorf("[TagsByResource] TagResourceModel.FindTagsByBizIDAndTargetID err: %v req: %+v", err, in)
		return nil, err
	}
	if len(trs) == 0 {
		return &pb.TagsByResourceResponse{}, nil
	}

	tagIds := make([]int64, len(trs))
	for i, tr := range trs {
		tagIds[i] = tr.TagID
	}

	tags, err := l.svcCtx.TagModel.FindByIds(l.ctx, tagIds)
	if err != nil {
		l.Logger.Errorf("[TagsByResource] TagModel.FindByIds err: %v tagIds: %v", err, tagIds)
		return nil, err
	}

	countMap, err := l.svcCtx.TagResourceModel.CountByTagIDs(l.ctx, tagIds)
	if err != nil {
		l.Logger.Errorf("[TagsByResource] TagResourceModel.CountByTagIDs err: %v", err)
	}

	return &pb.TagsByResourceResponse{
		Items: buildTagItems(tags, countMap),
	}, nil
}
