package qalogic

import (
	"context"

	"rpc-content/content"
	"rpc-content/internal/svc"
	"rpc-content/pkg/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type AnswerDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAnswerDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AnswerDeleteLogic {
	return &AnswerDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AnswerDeleteLogic) AnswerDelete(in *content.AnswerDeleteRequest) (resp *content.AnswerDeleteResponse, err error) {
	resp = new(content.AnswerDeleteResponse)

	if in.UserId <= 0 {
		resp.Code = int64(code.QAUserIdInvalid.Code())
		resp.Msg = code.QAUserIdInvalid.Message()
		return resp, nil
	}
	if in.AnswerId == 0 {
		resp.Code = int64(code.AnswerNotFound.Code())
		resp.Msg = code.AnswerNotFound.Message()
		return resp, nil
	}

	a, err := l.svcCtx.AnswerModel.FindOne(l.ctx, in.AnswerId)
	if err != nil {
		resp.Code = int64(code.ServerErr.Code())
		resp.Msg = err.Error()
		return resp, nil
	}
	if a == nil || a.Status == 1 {
		resp.Code = int64(code.AnswerNotFound.Code())
		resp.Msg = code.AnswerNotFound.Message()
		return resp, nil
	}
	if a.AuthorID != in.UserId {
		resp.Code = int64(code.NotAnswerAuthor.Code())
		resp.Msg = code.NotAnswerAuthor.Message()
		return resp, nil
	}

	if err := l.svcCtx.AnswerModel.UpdateFields(l.ctx, in.AnswerId, map[string]interface{}{"status": 1}); err != nil {
		resp.Code = int64(code.ServerErr.Code())
		resp.Msg = err.Error()
		return resp, nil
	}
	_ = l.svcCtx.QuestionModel.DecrAnswerNum(l.ctx, a.QuestionID)

	return resp, nil
}
