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

func (l *QuestionDetailLogic) QuestionDetail(in *content.QuestionDetailRequest) (resp *content.QuestionDetailResponse, err error) {
	resp = new(content.QuestionDetailResponse)
	resp.Data = new(content.QuestionItem)

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

	resp.Data.Id = q.ID
	resp.Data.Title = q.Title
	resp.Data.Content = q.Content
	resp.Data.AuthorId = q.AuthorID
	resp.Data.AnswerNum = int64(q.AnswerNum)
	resp.Data.ViewNum = int64(q.ViewNum)
	resp.Data.TagIds = q.TagIds
	resp.Data.CreateTime = q.CreateTime.Unix()

	return resp, nil
}
