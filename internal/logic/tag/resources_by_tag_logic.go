package taglogic

import (
	"context"
	"math"

	"rpc-content/content"
	model "rpc-content/internal/model/tag"
	"rpc-content/internal/svc"
	types "rpc-content/internal/types/tag"
	"rpc-content/pkg/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResourcesByTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResourcesByTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourcesByTagLogic {
	return &ResourcesByTagLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResourcesByTagLogic) ResourcesByTag(in *content.ResourcesByTagRequest) (resp *content.ResourcesByTagResponse, err error) {
	resp = new(content.ResourcesByTagResponse)
	resp.Data = new(content.ResourcesByTagData)
	resp.Data.Items = make([]*content.ResourceItem, 0)

	if in.TagId == 0 {
		resp.Data = nil
		resp.Code = int64(code.TagIdEmpty.Code())
		resp.Msg = code.TagIdEmpty.Message()
		return resp, nil
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.Cursor == 0 {
		in.Cursor = math.MaxInt64
	}

	var (
		trs   []*model.TagResource
		isEnd bool
	)

	if in.BizId != "" {
		trs, err = l.svcCtx.TagResourceModel.FindResourcesByTagIDAndBizID(l.ctx, in.TagId, in.BizId, in.Cursor, in.PageSize+1)
	} else {
		trs, err = l.svcCtx.TagResourceModel.FindResourcesByTagID(l.ctx, in.TagId, in.Cursor, in.PageSize+1)
	}
	if err != nil {
		resp.Data = nil
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}

	if len(trs) > int(in.PageSize) {
		trs = trs[:in.PageSize]
	} else {
		isEnd = true
	}
	if len(trs) == 0 {
		resp.Data.IsEnd = true
		return resp, nil
	}

	items := make([]*content.ResourceItem, 0, len(trs))
	for _, tr := range trs {
		items = append(items, &content.ResourceItem{
			TargetId:   tr.TargetID,
			BizId:      tr.BizID,
			CreateTime: tr.CreateTime.Unix(),
		})
	}

	cursor := trs[len(trs)-1].ID
	resp.Data.Items = items
	resp.Data.Cursor = cursor
	resp.Data.IsEnd = isEnd
	return resp, nil
}
