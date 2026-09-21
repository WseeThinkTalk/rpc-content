package qalogic

import (
	"context"

	"rpc-content/content"
	"rpc-content/internal/svc"
	types "rpc-content/internal/types/qa"
	"rpc-content/pkg/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type AnswerListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAnswerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AnswerListLogic {
	return &AnswerListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AnswerListLogic) AnswerList(in *content.AnswerListRequest) (resp *content.AnswerListResponse, err error) {
	resp = new(content.AnswerListResponse)
	resp.Code = 200
	resp.Msg = "success"
	resp.Data = new(content.AnswerListData)
	resp.Data.Items = make([]*content.AnswerItem, 0)

	if in.QuestionId == 0 {
		resp.Code = int64(code.QuestionIdEmpty.Code())
		resp.Msg = code.QuestionIdEmpty.Message()
		return resp, nil
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}

	answers, err := l.svcCtx.AnswerModel.FindByQuestionId(l.ctx, in.QuestionId, in.Cursor, in.PageSize+1)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}

	var isEnd bool
	if len(answers) > int(in.PageSize) {
		answers = answers[:in.PageSize]
	} else {
		isEnd = true
	}
	if len(answers) == 0 {
		resp.Data.IsEnd = true
		return resp, nil
	}

	items := make([]*content.AnswerItem, 0, len(answers))
	for _, a := range answers {
		items = append(items, &content.AnswerItem{
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
	resp.Data.Items = items
	resp.Data.Cursor = cursor
	resp.Data.IsEnd = isEnd
	return resp, nil
}
