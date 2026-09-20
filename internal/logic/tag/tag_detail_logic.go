package taglogic

import (
	"context"

	"rpc-content/content"
	"rpc-content/internal/svc"
	"rpc-content/pkg/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTagDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagDetailLogic {
	return &TagDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TagDetailLogic) TagDetail(in *content.TagDetailRequest) (*content.TagDetailResponse, error) {
	if in.TagId == 0 {
		return nil, code.TagIdEmpty
	}

	tag, err := l.svcCtx.TagModel.FindOne(l.ctx, in.TagId)
	if err != nil {
		l.Logger.Errorf("[TagDetail] TagModel.FindOne err: %v tagId: %d", err, in.TagId)
		return nil, err
	}
	if tag == nil {
		return nil, code.TagNotFound
	}

	resourceCount, err := l.svcCtx.TagResourceModel.CountByTagID(l.ctx, in.TagId)
	if err != nil {
		l.Logger.Errorf("[TagDetail] TagResourceModel.CountByTagID err: %v tagId: %d", err, in.TagId)
		return nil, err
	}

	return &content.TagDetailResponse{
		Code: 200,
		Msg:  "success",
		Data: &content.TagItem{
			TagId:         tag.ID,
			TagName:       tag.TagName,
			TagDesc:       tag.TagDesc,
			ResourceCount: resourceCount,
			CreateTime:    tag.CreateTime.Unix(),
		},
	}, nil
}
