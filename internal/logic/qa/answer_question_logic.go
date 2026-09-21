package qalogic

import (
	"context"
	"time"

	"rpc-content/content"
	model "rpc-content/internal/model/qa"
	"rpc-content/internal/svc"
	"rpc-content/pkg/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type AnswerQuestionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAnswerQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AnswerQuestionLogic {
	return &AnswerQuestionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AnswerQuestionLogic) AnswerQuestion(in *content.AnswerQuestionRequest) (resp *content.AnswerQuestionResponse, err error) {
	resp = new(content.AnswerQuestionResponse)
	resp.Data = new(content.AnswerQuestionData)

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
	if len(in.Content) == 0 {
		resp.Code = int64(code.ContentEmpty.Code())
		resp.Msg = code.ContentEmpty.Message()
		return resp, nil
	}

	q, err := l.svcCtx.QuestionModel.FindOne(l.ctx, in.QuestionId)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}
	if q == nil || q.Status == 1 {
		resp.Code = int64(code.QuestionNotFound.Code())
		resp.Msg = code.QuestionNotFound.Message()
		return resp, nil
	}

	ans := &model.Answer{
		QuestionID: in.QuestionId,
		AuthorID:   in.UserId,
		Content:    in.Content,
		Status:     0,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	if err := l.svcCtx.AnswerModel.Insert(l.ctx, ans); err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}

	_ = l.svcCtx.QuestionModel.IncrAnswerNum(l.ctx, in.QuestionId)

	resp.Data.AnswerId = ans.ID
	return resp, nil
}
