package qalogic

import (
	"context"

	"rpc-content/content"
	"rpc-content/internal/svc"
	"rpc-content/pkg/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type AcceptAnswerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAcceptAnswerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AcceptAnswerLogic {
	return &AcceptAnswerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AcceptAnswerLogic) AcceptAnswer(in *content.AcceptAnswerRequest) (resp *content.AcceptAnswerResponse, err error) {
	resp = new(content.AcceptAnswerResponse)
	resp.Code = 200
	resp.Msg = "success"

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
	if in.AnswerId == 0 {
		resp.Code = int64(code.AnswerNotFound.Code())
		resp.Msg = code.AnswerNotFound.Message()
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

	accepted, _ := l.svcCtx.AnswerModel.FindAcceptedByQuestionId(l.ctx, in.QuestionId)
	if accepted != nil {
		resp.Code = int64(code.AlreadyAccepted.Code())
		resp.Msg = code.AlreadyAccepted.Message()
		return resp, nil
	}

	a, err := l.svcCtx.AnswerModel.FindOne(l.ctx, in.AnswerId)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}
	if a == nil || a.QuestionID != in.QuestionId {
		resp.Code = int64(code.AnswerNotFound.Code())
		resp.Msg = code.AnswerNotFound.Message()
		return resp, nil
	}

	if err := l.svcCtx.AnswerModel.UpdateFields(l.ctx, in.AnswerId, map[string]interface{}{"is_accepted": 1}); err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}

	return resp, nil
}
