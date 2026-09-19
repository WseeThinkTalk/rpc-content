package logic

import (
	"context"

	"rpc-content/qa/code"
	"rpc-content/qa/internal/svc"
	"rpc-content/qa/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AcceptAnswerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAcceptAnswerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AcceptAnswerLogic {
	return &AcceptAnswerLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AcceptAnswerLogic) AcceptAnswer(in *pb.AcceptAnswerRequest) (*pb.AcceptAnswerResponse, error) {
	if in.UserId <= 0 {
		return nil, code.UserIdInvalid
	}
	if in.QuestionId == 0 {
		return nil, code.QuestionIdEmpty
	}
	if in.AnswerId == 0 {
		return nil, code.AnswerNotFound
	}

	q, err := l.svcCtx.QuestionModel.FindOne(l.ctx, in.QuestionId)
	if err != nil {
		l.Errorf("[AcceptAnswer] FindOne question err: %v id: %d", err, in.QuestionId)
		return nil, err
	}
	if q == nil {
		return nil, code.QuestionNotFound
	}
	if q.AuthorID != in.UserId {
		return nil, code.NotQuestionAuthor
	}

	accepted, _ := l.svcCtx.AnswerModel.FindAcceptedByQuestionId(l.ctx, in.QuestionId)
	if accepted != nil {
		return nil, code.AlreadyAccepted
	}

	a, err := l.svcCtx.AnswerModel.FindOne(l.ctx, in.AnswerId)
	if err != nil {
		l.Errorf("[AcceptAnswer] FindOne answer err: %v id: %d", err, in.AnswerId)
		return nil, err
	}
	if a == nil || a.QuestionID != in.QuestionId {
		return nil, code.AnswerNotFound
	}

	if err := l.svcCtx.AnswerModel.UpdateFields(l.ctx, in.AnswerId, map[string]interface{}{"is_accepted": 1}); err != nil {
		l.Errorf("[AcceptAnswer] UpdateFields err: %v id: %d", err, in.AnswerId)
		return nil, err
	}

	return &pb.AcceptAnswerResponse{}, nil
}
