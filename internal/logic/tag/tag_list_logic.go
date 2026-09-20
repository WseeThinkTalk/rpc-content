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

func (l *TagListLogic) TagList(in *content.TagListRequest) (*content.TagListResponse, error) {
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.Cursor == 0 {
		in.Cursor = math.MaxInt64
	}

	tags, err := l.svcCtx.TagModel.FindByCursor(l.ctx, in.Cursor, in.PageSize+1)
	if err != nil {
		l.Logger.Errorf("[TagList] TagModel.FindByCursor err: %v cursor: %d", err, in.Cursor)
		return nil, err
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
		return &content.TagListResponse{
			Code: 200,
			Msg:  "success",
			Data: &content.TagListData{
				Items: []*content.TagItem{},
				IsEnd: true,
			},
		}, nil
	}

	tagIds := make([]int64, len(tags))
	for i, t := range tags {
		tagIds[i] = t.ID
	}
	countMap, err := l.svcCtx.TagResourceModel.CountByTagIDs(l.ctx, tagIds)
	if err != nil {
		l.Logger.Errorf("[TagList] TagResourceModel.CountByTagIDs err: %v", err)
	}

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
	return &content.TagListResponse{
		Code: 200,
		Msg:  "success",
		Data: &content.TagListData{
			Items:  items,
			Cursor: cursor,
			IsEnd:  isEnd,
		},
	}, nil
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
