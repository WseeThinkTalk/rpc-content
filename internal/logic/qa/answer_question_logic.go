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

func (l *AnswerQuestionLogic) AnswerQuestion(in *content.AnswerQuestionRequest) (*content.AnswerQuestionResponse, error) {
	if in.UserId <= 0 {
		return nil, code.QAUserIdInvalid
	}
	if in.QuestionId == 0 {
		return nil, code.QuestionIdEmpty
	}
	if len(in.Content) == 0 {
		return nil, code.ContentEmpty
	}

	q, err := l.svcCtx.QuestionModel.FindOne(l.ctx, in.QuestionId)
	if err != nil {
		l.Errorf("[AnswerQuestion] FindOne question err: %v id: %d", err, in.QuestionId)
		return nil, err
	}
	if q == nil || q.Status == 1 {
		return nil, code.QuestionNotFound
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
		l.Errorf("[AnswerQuestion] Insert err: %v req: %+v", err, in)
		return nil, err
	}

	_ = l.svcCtx.QuestionModel.IncrAnswerNum(l.ctx, in.QuestionId)

	return &content.AnswerQuestionResponse{
		Code: 200,
		Msg:  "success",
		Data: &content.AnswerQuestionData{
			AnswerId: ans.ID,
		},
	}, nil
}
