package taglogic

import (
	"context"
	"math"

	"rpc-content/content"
	model "rpc-content/internal/model/tag"
	"rpc-content/internal/svc"
	types "rpc-content/internal/types/tag"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTagListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagListLogic {
	return &TagListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TagListLogic) TagList(in *content.TagListRequest) (resp *content.TagListResponse, err error) {
	resp = new(content.TagListResponse)
	resp.Code = 200
	resp.Msg = "success"
	resp.Data = new(content.TagListData)
	resp.Data.Items = make([]*content.TagItem, 0)

	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.Cursor == 0 {
		in.Cursor = math.MaxInt64
	}

	tags, err := l.svcCtx.TagModel.FindByCursor(l.ctx, in.Cursor, in.PageSize+1)
	if err != nil {
		resp.Data = nil
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, nil
	}

	var (
		isEnd  bool
		cursor int64
		items  []*content.TagItem
	)
	if len(tags) > int(in.PageSize) {
		tags = tags[:in.PageSize]
	} else {
		isEnd = true
	}
	if len(tags) == 0 {
		resp.Data.IsEnd = true
		return resp, nil
	}

	tagIds := make([]int64, len(tags))
	for i, t := range tags {
		tagIds[i] = t.ID
	}
	countMap, _ := l.svcCtx.TagResourceModel.CountByTagIDs(l.ctx, tagIds)

	for _, t := range tags {
		items = append(items, &content.TagItem{
			TagId:         t.ID,
			TagName:       t.TagName,
			TagDesc:       t.TagDesc,
			ResourceCount: countMap[t.ID],
			CreateTime:    t.CreateTime.Unix(),
		})
	}

	cursor = tags[len(tags)-1].ID
	resp.Data.Items = items
	resp.Data.Cursor = cursor
	resp.Data.IsEnd = isEnd
	return resp, nil
}

func buildTagItems(tags []*model.Tag, countMap map[int64]int64) []*content.TagItem {
	items := make([]*content.TagItem, 0, len(tags))
	for _, t := range tags {
		items = append(items, &content.TagItem{
			TagId:         t.ID,
			TagName:       t.TagName,
			TagDesc:       t.TagDesc,
			ResourceCount: countMap[t.ID],
			CreateTime:    t.CreateTime.Unix(),
		})
	}
	return items
}
