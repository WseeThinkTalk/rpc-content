package taglogic

import (
	"context"

	"rpc-content/content"
	"rpc-content/internal/svc"
	types "rpc-content/internal/types/tag"

	"github.com/zeromicro/go-zero/core/logx"
)

type HotTagsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHotTagsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HotTagsLogic {
	return &HotTagsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HotTagsLogic) HotTags(in *content.HotTagsRequest) (resp *content.HotTagsResponse, err error) {
	resp = new(content.HotTagsResponse)
	resp.Data = make([]*content.TagItem, 0)

	limit := int(in.Limit)
	if limit <= 0 || limit > types.HotTagsMaxCount {
		limit = types.HotTagsMaxCount
	}

	tagIds, err := l.svcCtx.TagResourceModel.FindHotTagIDs(l.ctx, limit)
	if err != nil {
		l.Logger.Errorf("[HotTags] TagResourceModel.FindHotTagIDs err: %v", err)
		return nil, err
	}
	if len(tagIds) == 0 {
		return resp, nil
	}

	tags, err := l.svcCtx.TagModel.FindByIds(l.ctx, tagIds)
	if err != nil {
		l.Logger.Errorf("[HotTags] TagModel.FindByIds err: %v tagIds: %v", err, tagIds)
		return nil, err
	}

	countMap, err := l.svcCtx.TagResourceModel.CountByTagIDs(l.ctx, tagIds)
	if err != nil {
		l.Logger.Errorf("[HotTags] TagResourceModel.CountByTagIDs err: %v", err)
	}

	resp.Data = buildTagItems(tags, countMap)
	return resp, nil
}
