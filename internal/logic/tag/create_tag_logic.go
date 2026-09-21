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

type CreateTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTagLogic {
	return &CreateTagLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateTagLogic) CreateTag(in *content.CreateTagRequest) (resp *content.CreateTagResponse, err error) {
	resp = new(content.CreateTagResponse)
	resp.Code = 200
	resp.Msg = "success"
	resp.Data = new(content.CreateTagData)

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

	exist, err := l.svcCtx.TagModel.FindByName(l.ctx, in.TagName)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}
	if exist != nil {
		resp.Code = int64(code.TagNameExists.Code())
		resp.Msg = code.TagNameExists.Message()
		return resp, nil
	}

	tag := &model.Tag{
		TagName:    in.TagName,
		TagDesc:    in.TagDesc,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	if err := l.svcCtx.TagModel.Insert(l.ctx, tag); err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}

	resp.Data.TagId = tag.ID
	return resp, nil
}
