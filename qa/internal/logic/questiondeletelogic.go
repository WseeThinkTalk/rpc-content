package logic

import (
	"context"

	"rpc-content/qa/code"
	"rpc-content/qa/internal/svc"
	"rpc-content/qa/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuestionDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQuestionDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuestionDeleteLogic {
	return &QuestionDeleteLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *QuestionDeleteLogic) QuestionDelete(in *pb.QuestionDeleteRequest) (*pb.QuestionDeleteResponse, error) {
	if in.UserId <= 0 {
		return nil, code.UserIdInvalid
	}
	if in.QuestionId == 0 {
		return nil, code.QuestionIdEmpty
	}

	q, err := l.svcCtx.QuestionModel.FindOne(l.ctx, in.QuestionId)
	if err != nil {
		l.Errorf("[QuestionDelete] FindOne err: %v id: %d", err, in.QuestionId)
		return nil, err
	}
	if q == nil {
		return nil, code.QuestionNotFound
	}
	if q.AuthorID != in.UserId {
		return nil, code.NotQuestionAuthor
	}

	if err := l.svcCtx.QuestionModel.UpdateFields(l.ctx, in.QuestionId, map[string]interface{}{"status": 1}); err != nil {
		l.Errorf("[QuestionDelete] UpdateFields err: %v id: %d", err, in.QuestionId)
		return nil, err
	}

	return &pb.QuestionDeleteResponse{}, nil
}
