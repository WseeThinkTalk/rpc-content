package qalogic

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"rpc-content/content"
	model "rpc-content/internal/model/qa"
	"rpc-content/internal/svc"
	types "rpc-content/internal/types/qa"
	"rpc-content/pkg/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishQuestionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishQuestionLogic {
	return &PublishQuestionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishQuestionLogic) PublishQuestion(in *content.PublishQuestionRequest) (resp *content.PublishQuestionResponse, err error) {
	resp = new(content.PublishQuestionResponse)
	resp.Code = 200
	resp.Msg = "success"
	resp.Data = new(content.PublishQuestionData)

	if in.UserId <= 0 {
		resp.Code = int64(code.QAUserIdInvalid.Code())
		resp.Msg = code.QAUserIdInvalid.Message()
		return resp, nil
	}
	if len(in.Title) == 0 {
		resp.Code = int64(code.TitleEmpty.Code())
		resp.Msg = code.TitleEmpty.Message()
		return resp, nil
	}
	if len(in.Content) == 0 {
		resp.Code = int64(code.ContentEmpty.Code())
		resp.Msg = code.ContentEmpty.Message()
		return resp, nil
	}

	q := &model.Question{
		Title:      in.Title,
		Content:    in.Content,
		AuthorID:   in.UserId,
		TagIds:     in.TagIds,
		Status:     types.QuestionStatusNormal,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	if err := l.svcCtx.QuestionModel.Insert(l.ctx, q); err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}

	key := questionsKey(in.UserId, types.SortPublishTime)
	b, _ := l.svcCtx.BizRedis.ExistsCtx(l.ctx, key)
	if b {
		_, _ = l.svcCtx.BizRedis.ZaddCtx(l.ctx, key, time.Now().Unix(), strconv.FormatInt(q.ID, 10))
		_ = l.svcCtx.BizRedis.ExpireCtx(l.ctx, key, types.CacheExpireTime)
	}

	resp.Data.QuestionId = q.ID
	return resp, nil
}

func questionsKey(uid int64, sortType int32) string {
	return fmt.Sprintf("biz#questions#%d#%d", uid, sortType)
}
