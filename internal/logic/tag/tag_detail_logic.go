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

func (l *TagDetailLogic) TagDetail(in *content.TagDetailRequest) (resp *content.TagDetailResponse, err error) {
	resp = new(content.TagDetailResponse)
	resp.Data = new(content.TagItem)

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

	resp.Data.TagId = tag.ID
	resp.Data.TagName = tag.TagName
	resp.Data.TagDesc = tag.TagDesc
	resp.Data.ResourceCount = resourceCount
	resp.Data.CreateTime = tag.CreateTime.Unix()

	return resp, nil
}
