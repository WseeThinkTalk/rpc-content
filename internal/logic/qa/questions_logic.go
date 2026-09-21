package qalogic

import (
	"context"
	"time"

	"rpc-content/content"
	"rpc-content/internal/svc"
	types "rpc-content/internal/types/qa"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuestionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQuestionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuestionsLogic {
	return &QuestionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QuestionsLogic) Questions(in *content.QuestionsRequest) (resp *content.QuestionsResponse, err error) {
	resp = new(content.QuestionsResponse)
	resp.Code = 200
	resp.Msg = "success"
	resp.Data = new(content.QuestionsData)
	resp.Data.Items = make([]*content.QuestionItem, 0)

	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.SortType == 1 && in.Cursor == 0 {
		in.Cursor = 999999
	}
	if in.SortType == 0 && in.Cursor == 0 {
		in.Cursor = time.Now().UnixNano() / 1e6
	}

	questions, err := l.svcCtx.QuestionModel.QuestionsByUserId(l.ctx, in.UserId, int(in.SortType), in.Cursor, in.PageSize+1)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}

	var isEnd bool
	if len(questions) > int(in.PageSize) {
		questions = questions[:in.PageSize]
	} else {
		isEnd = true
	}
	if len(questions) == 0 {
		resp.Data.IsEnd = true
		return resp, nil
	}

	items := make([]*content.QuestionItem, 0, len(questions))
	for _, q := range questions {
		items = append(items, &content.QuestionItem{
			Id:         q.ID,
			Title:      q.Title,
			Content:    q.Content,
			AuthorId:   q.AuthorID,
			AnswerNum:  int64(q.AnswerNum),
			ViewNum:    int64(q.ViewNum),
			TagIds:     q.TagIds,
			CreateTime: q.CreateTime.Unix(),
		})
	}

	var cursor int64
	last := questions[len(questions)-1]
	if in.SortType == types.SortHot {
		cursor = int64(last.AnswerNum)
	} else {
		cursor = last.ID
	}

	resp.Data.Items = items
	resp.Data.Cursor = cursor
	resp.Data.IsEnd = isEnd
	return resp, nil
}
