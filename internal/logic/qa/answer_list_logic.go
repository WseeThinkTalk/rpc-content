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

func (l *AnswerListLogic) AnswerList(in *content.AnswerListRequest) (*content.AnswerListResponse, error) {
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
		return &content.AnswerListResponse{
			Code: 200,
			Msg:  "success",
			Data: &content.AnswerListData{
				Items: []*content.AnswerItem{},
				IsEnd: true,
			},
		}, nil
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
	return &content.AnswerListResponse{
		Code: 200,
		Msg:  "success",
		Data: &content.AnswerListData{
			Items:  items,
			Cursor: cursor,
			IsEnd:  isEnd,
		},
	}, nil
}
