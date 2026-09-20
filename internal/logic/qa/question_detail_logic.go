package qalogic

import (
	"context"

	"rpc-content/content"
	"rpc-content/internal/svc"
	"rpc-content/pkg/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuestionDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQuestionDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuestionDetailLogic {
	return &QuestionDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QuestionDetailLogic) QuestionDetail(in *content.QuestionDetailRequest) (*content.QuestionDetailResponse, error) {
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

	item := &content.QuestionItem{
		Id:         q.ID,
		Title:      q.Title,
		Content:    q.Content,
		AuthorId:   q.AuthorID,
		AnswerNum:  int64(q.AnswerNum),
		ViewNum:    int64(q.ViewNum),
		TagIds:     q.TagIds,
		CreateTime: q.CreateTime.Unix(),
	}

	return &content.QuestionDetailResponse{
		Code: 200,
		Msg:  "success",
		Data: item,
	}, nil
}
