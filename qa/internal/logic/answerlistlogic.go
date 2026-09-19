package logic

import (
	"context"

	"rpc-content/qa/code"
	"rpc-content/qa/internal/svc"
	"rpc-content/qa/internal/types"
	"rpc-content/qa/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AnswerListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAnswerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AnswerListLogic {
	return &AnswerListLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AnswerListLogic) AnswerList(in *pb.AnswerListRequest) (*pb.AnswerListResponse, error) {
	if in.QuestionId == 0 {
		return nil, code.QuestionIdEmpty
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}

	answers, err := l.svcCtx.AnswerModel.FindByQuestionId(l.ctx, in.QuestionId, in.Cursor, in.PageSize+1)
	if err != nil {
		l.Errorf("[AnswerList] FindByQuestionId err: %v questionId: %d", err, in.QuestionId)
		return nil, err
	}

	var isEnd bool
	if len(answers) > int(in.PageSize) {
		answers = answers[:in.PageSize]
	} else {
		isEnd = true
	}
	if len(answers) == 0 {
		return &pb.AnswerListResponse{IsEnd: true}, nil
	}

	items := make([]*pb.AnswerItem, 0, len(answers))
	for _, a := range answers {
		items = append(items, &pb.AnswerItem{
			Id:         a.ID,
			QuestionId: a.QuestionID,
			AuthorId:   a.AuthorID,
			Content:    a.Content,
			IsAccepted: a.IsAccepted == 1,
			LikeNum:    int64(a.LikeNum),
			ReplyNum:   int64(a.ReplyNum),
			CreateTime: a.CreateTime.Unix(),
		})
	}

	cursor := answers[len(answers)-1].ID
	return &pb.AnswerListResponse{
		Items:  items,
		Cursor: cursor,
		IsEnd:  isEnd,
	}, nil
}
