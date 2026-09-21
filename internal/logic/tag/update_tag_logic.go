package taglogic

import (
	"context"

	"rpc-content/content"
	"rpc-content/internal/svc"
	"rpc-content/pkg/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTagLogic {
	return &UpdateTagLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateTagLogic) UpdateTag(in *content.UpdateTagRequest) (resp *content.UpdateTagResponse, err error) {
	resp = new(content.UpdateTagResponse)

	if in.TagId == 0 {
		resp.Code = int64(code.TagIdEmpty.Code())
		resp.Msg = code.TagIdEmpty.Message()
		return resp, nil
	}
	if in.TagName == "" {
		resp.Code = int64(code.TagNameEmpty.Code())
		resp.Msg = code.TagNameEmpty.Message()
		return resp, nil
	}
	if len(in.TagName) > 32 {
		resp.Code = int64(code.TagNameTooLong.Code())
		resp.Msg = code.TagNameTooLong.Message()
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

	if in.TagName != tag.TagName {
		exist, err := l.svcCtx.TagModel.FindByName(l.ctx, in.TagName)
		if err != nil {
			resp.Code = int64(code.ServerErr.Code())
			resp.Msg = err.Error()
			return resp, nil
		}
		if exist != nil {
			resp.Code = int64(code.TagNameExists.Code())
			resp.Msg = code.TagNameExists.Message()
			return resp, nil
		}
	}

	err = l.svcCtx.TagModel.UpdateFields(l.ctx, in.TagId, map[string]interface{}{
		"tag_name": in.TagName,
		"tag_desc": in.TagDesc,
	})
	if err != nil {
		resp.Code = int64(code.ServerErr.Code())
		resp.Msg = err.Error()
		return resp, nil
	}

	return resp, nil
}
