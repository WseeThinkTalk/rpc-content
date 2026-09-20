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

func (l *UpdateTagLogic) UpdateTag(in *content.UpdateTagRequest) (*content.UpdateTagResponse, error) {
	if in.TagId == 0 {
		return nil, code.TagIdEmpty
	}
	if in.TagName == "" {
		return nil, code.TagNameEmpty
	}
	if len(in.TagName) > 32 {
		return nil, code.TagNameTooLong
	}

	tag, err := l.svcCtx.TagModel.FindOne(l.ctx, in.TagId)
	if err != nil {
		l.Logger.Errorf("[UpdateTag] TagModel.FindOne err: %v tagId: %d", err, in.TagId)
		return nil, err
	}
	if tag == nil {
		return nil, code.TagNotFound
	}

	if in.TagName != tag.TagName {
		exist, err := l.svcCtx.TagModel.FindByName(l.ctx, in.TagName)
		if err != nil {
			l.Logger.Errorf("[UpdateTag] TagModel.FindByName err: %v tagName: %s", err, in.TagName)
			return nil, err
		}
		if exist != nil {
			return nil, code.TagNameExists
		}
	}

	err = l.svcCtx.TagModel.UpdateFields(l.ctx, in.TagId, map[string]interface{}{
		"tag_name": in.TagName,
		"tag_desc": in.TagDesc,
	})
	if err != nil {
		l.Logger.Errorf("[UpdateTag] TagModel.UpdateFields err: %v tagId: %d", err, in.TagId)
		return nil, err
	}

	return &content.UpdateTagResponse{
		Code: 200,
		Msg:  "success",
	}, nil
}
