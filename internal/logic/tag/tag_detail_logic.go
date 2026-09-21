package taglogic

import (
	"context"

	"rpc-content/content"
	"rpc-content/internal/svc"
	"rpc-content/pkg/code"

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

func (l *TagDetailLogic) TagDetail(in *content.TagDetailRequest) (resp *content.TagDetailResponse, err error) {
	resp = new(content.TagDetailResponse)
	resp.Data = new(content.TagItem)

	if in.TagId == 0 {
		resp.Code = int64(code.TagIdEmpty.Code())
		resp.Msg = code.TagIdEmpty.Message()
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

	resourceCount, err := l.svcCtx.TagResourceModel.CountByTagID(l.ctx, in.TagId)
	if err != nil {
		resp.Code = int64(code.ServerErr.Code())
		resp.Msg = err.Error()
		return resp, nil
	}

	resp.Data.TagId = tag.ID
	resp.Data.TagName = tag.TagName
	resp.Data.TagDesc = tag.TagDesc
	resp.Data.ResourceCount = resourceCount
	resp.Data.CreateTime = tag.CreateTime.Unix()

	return resp, nil
}
