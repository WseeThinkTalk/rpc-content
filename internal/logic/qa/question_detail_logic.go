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
		resp.Data = nil
		resp.Code = int64(code.QuestionIdEmpty.Code())
		resp.Msg = code.QuestionIdEmpty.Message()
		return resp, nil
	}

	q, err := l.svcCtx.QuestionModel.FindOne(l.ctx, in.QuestionId)
	if err != nil {
		resp.Data = nil
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}
	if q == nil || q.Status == 1 {
		resp.Data = nil
		resp.Code = int64(code.QuestionNotFound.Code())
		resp.Msg = code.QuestionNotFound.Message()
		return resp, nil
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
