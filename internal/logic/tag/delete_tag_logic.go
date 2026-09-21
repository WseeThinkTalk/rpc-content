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

func (l *DeleteTagLogic) DeleteTag(in *content.DeleteTagRequest) (resp *content.DeleteTagResponse, err error) {
	resp = new(content.DeleteTagResponse)

	if in.TagId == 0 {
		resp.Code = int64(code.TagIdEmpty.Code())
		resp.Msg = code.TagIdEmpty.Message()
		return resp, nil
	}

	tag, err := l.svcCtx.TagModel.FindOne(l.ctx, in.TagId)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}
	if tag == nil {
		resp.Code = int64(code.TagNotFound.Code())
		resp.Msg = code.TagNotFound.Message()
		return resp, nil
	}

	err = l.svcCtx.TagModel.Delete(l.ctx, in.TagId)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}

	return resp, nil
}
