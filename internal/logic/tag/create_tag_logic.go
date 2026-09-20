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

func (l *CreateTagLogic) CreateTag(in *content.CreateTagRequest) (*content.CreateTagResponse, error) {
	if in.TagName == "" {
		return nil, code.TagNameEmpty
	}
	if len(in.TagName) > 32 {
		return nil, code.TagNameTooLong
	}

	exist, err := l.svcCtx.TagModel.FindByName(l.ctx, in.TagName)
	if err != nil {
		l.Logger.Errorf("[CreateTag] TagModel.FindByName err: %v tagName: %s", err, in.TagName)
		return nil, err
	}
	if exist != nil {
		return nil, code.TagNameExists
	}

	tag := &model.Tag{
		TagName:    in.TagName,
		TagDesc:    in.TagDesc,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	if err := l.svcCtx.TagModel.Insert(l.ctx, tag); err != nil {
		l.Logger.Errorf("[CreateTag] TagModel.Insert err: %v tag: %+v", err, tag)
		return nil, err
	}

	return &content.CreateTagResponse{
		Code: 200,
		Msg:  "success",
		Data: &content.CreateTagData{
			TagId: tag.ID,
		},
	}, nil
}
