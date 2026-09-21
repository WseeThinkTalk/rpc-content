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
		resp.Code = int64(code.BizIdEmpty.Code())
		resp.Msg = code.BizIdEmpty.Message()
		return resp, nil
	}
	if in.TargetId == 0 {
		resp.Code = int64(code.TargetIdEmpty.Code())
		resp.Msg = code.TargetIdEmpty.Message()
		return resp, nil
	}

	trs, err := l.svcCtx.TagResourceModel.FindTagsByBizIDAndTargetID(l.ctx, in.BizId, in.TargetId)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
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
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}

	countMap, _ := l.svcCtx.TagResourceModel.CountByTagIDs(l.ctx, tagIds)

	resp.Data = buildTagItems(tags, countMap)
	return resp, nil
}
