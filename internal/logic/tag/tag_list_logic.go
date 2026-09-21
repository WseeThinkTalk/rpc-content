package taglogic

import (
	"rpc-content/pkg/code"
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
		resp.Code = int64(code.ServerErr.Code())
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

	// 提取标签ID列表
	tagIds := make([]int64, len(tags))
	for i, v := range tags {
		tagIds[i] = v.ID
	}
	countMap, _ := l.svcCtx.TagResourceModel.CountByTagIDs(l.ctx, tagIds)

	// 组装标签列表数据项
	for _, v := range tags {
		items = append(items, &content.TagItem{
			TagId:         v.ID,
			TagName:       v.TagName,
			TagDesc:       v.TagDesc,
			ResourceCount: countMap[v.ID],
			CreateTime:    v.CreateTime.Unix(),
		})
	}

	cursor = tags[len(tags)-1].ID
	resp.Data.Items = items
	resp.Data.Cursor = cursor
	resp.Data.IsEnd = isEnd
	return resp, nil
}

func buildTagItems(tags []*model.Tag, countMap map[int64]int64) []*content.TagItem {
	// 构建标签响应数据列表
	items := make([]*content.TagItem, 0, len(tags))
	for _, v := range tags {
		items = append(items, &content.TagItem{
			TagId:         v.ID,
			TagName:       v.TagName,
			TagDesc:       v.TagDesc,
			ResourceCount: countMap[v.ID],
			CreateTime:    v.CreateTime.Unix(),
		})
	}
	return items
}
