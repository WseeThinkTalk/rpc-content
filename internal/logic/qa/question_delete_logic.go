package qalogic

import (
	"context"

	"rpc-content/content"
	"rpc-content/internal/svc"
	"rpc-content/pkg/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuestionDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQuestionDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuestionDeleteLogic {
	return &QuestionDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QuestionDeleteLogic) QuestionDelete(in *content.QuestionDeleteRequest) (resp *content.QuestionDeleteResponse, err error) {
	resp = new(content.QuestionDeleteResponse)

	if in.UserId <= 0 {
		resp.Code = int64(code.QAUserIdInvalid.Code())
		resp.Msg = code.QAUserIdInvalid.Message()
		return resp, nil
	}
	if in.QuestionId == 0 {
		resp.Code = int64(code.QuestionIdEmpty.Code())
		resp.Msg = code.QuestionIdEmpty.Message()
		return resp, nil
	}

	q, err := l.svcCtx.QuestionModel.FindOne(l.ctx, in.QuestionId)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}
	if q == nil {
		resp.Code = int64(code.QuestionNotFound.Code())
		resp.Msg = code.QuestionNotFound.Message()
		return resp, nil
	}
	if q.AuthorID != in.UserId {
		resp.Code = int64(code.NotQuestionAuthor.Code())
		resp.Msg = code.NotQuestionAuthor.Message()
		return resp, nil
	}

	if err := l.svcCtx.QuestionModel.UpdateFields(l.ctx, in.QuestionId, map[string]interface{}{"status": 1}); err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}

	return resp, nil
}
