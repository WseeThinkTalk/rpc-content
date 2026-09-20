package taglogic

import (
	"context"

	"rpc-content/content"
	"rpc-content/internal/svc"
	"rpc-content/pkg/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTagLogic {
	return &DeleteTagLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteTagLogic) DeleteTag(in *content.DeleteTagRequest) (*content.DeleteTagResponse, error) {
	if in.TagId == 0 {
		return nil, code.TagIdEmpty
	}

	tag, err := l.svcCtx.TagModel.FindOne(l.ctx, in.TagId)
	if err != nil {
		l.Logger.Errorf("[DeleteTag] TagModel.FindOne err: %v tagId: %d", err, in.TagId)
		return nil, err
	}
	if tag == nil {
		return nil, code.TagNotFound
	}

	err = l.svcCtx.TagModel.Delete(l.ctx, in.TagId)
	if err != nil {
		l.Logger.Errorf("[DeleteTag] TagModel.Delete err: %v tagId: %d", err, in.TagId)
		return nil, err
	}

	return &content.DeleteTagResponse{
		Code: 200,
		Msg:  "success",
	}, nil
}
