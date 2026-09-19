package logic

import (
	"context"
	"time"

	"rpc-content/tag/code"
	"rpc-content/tag/internal/model"
	"rpc-content/tag/internal/svc"
	"rpc-content/tag/pb"

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

func (l *CreateTagLogic) CreateTag(in *pb.CreateTagRequest) (*pb.CreateTagResponse, error) {
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

	return &pb.CreateTagResponse{TagId: tag.ID}, nil
}
