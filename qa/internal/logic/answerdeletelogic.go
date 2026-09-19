package logic

import (
	"context"

	"rpc-content/qa/code"
	"rpc-content/qa/internal/svc"
	"rpc-content/qa/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AnswerDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAnswerDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AnswerDeleteLogic {
	return &AnswerDeleteLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AnswerDeleteLogic) AnswerDelete(in *pb.AnswerDeleteRequest) (*pb.AnswerDeleteResponse, error) {
	if in.UserId <= 0 {
		return nil, code.UserIdInvalid
	}
	if in.AnswerId == 0 {
		return nil, code.AnswerNotFound
	}

	a, err := l.svcCtx.AnswerModel.FindOne(l.ctx, in.AnswerId)
	if err != nil {
		l.Errorf("[AnswerDelete] FindOne err: %v id: %d", err, in.AnswerId)
		return nil, err
	}
	if a == nil || a.Status == 1 {
		return nil, code.AnswerNotFound
	}
	if a.AuthorID != in.UserId {
		return nil, code.NotAnswerAuthor
	}

	if err := l.svcCtx.AnswerModel.UpdateFields(l.ctx, in.AnswerId, map[string]interface{}{"status": 1}); err != nil {
		l.Errorf("[AnswerDelete] UpdateFields err: %v id: %d", err, in.AnswerId)
		return nil, err
	}
	_ = l.svcCtx.QuestionModel.DecrAnswerNum(l.ctx, a.QuestionID)

	return &pb.AnswerDeleteResponse{}, nil
}
