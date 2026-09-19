package logic

import (
	"context"

	"rpc-content/qa/code"
	"rpc-content/qa/internal/svc"
	"rpc-content/qa/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuestionDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQuestionDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuestionDetailLogic {
	return &QuestionDetailLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *QuestionDetailLogic) QuestionDetail(in *pb.QuestionDetailRequest) (*pb.QuestionDetailResponse, error) {
	if in.QuestionId == 0 {
		return nil, code.QuestionIdEmpty
	}

	q, err := l.svcCtx.QuestionModel.FindOne(l.ctx, in.QuestionId)
	if err != nil {
		l.Errorf("[QuestionDetail] FindOne err: %v id: %d", err, in.QuestionId)
		return nil, err
	}
	if q == nil || q.Status == 1 {
		return nil, code.QuestionNotFound
	}

	item := &pb.QuestionItem{
		Id:         q.ID,
		Title:      q.Title,
		Content:    q.Content,
		AuthorId:   q.AuthorID,
		AnswerNum:  int64(q.AnswerNum),
		ViewNum:    int64(q.ViewNum),
		TagIds:     q.TagIds,
		CreateTime: q.CreateTime.Unix(),
	}

	return &pb.QuestionDetailResponse{Question: item}, nil
}
