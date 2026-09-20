package taglogic

import (
	"context"

	"rpc-content/content"
	"rpc-content/internal/svc"
	"rpc-content/pkg/code"

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

func (l *TagsByResourceLogic) TagsByResource(in *content.TagsByResourceRequest) (resp *content.TagsByResourceResponse, err error) {
	resp = new(content.TagsByResourceResponse)
	resp.Code = 200
	resp.Msg = "success"
	resp.Data = make([]*content.TagItem, 0)

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
		return resp, nil
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

	resp.Data = buildTagItems(tags, countMap)
	return resp, nil
}
