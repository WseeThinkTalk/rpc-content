package articlelogic

import (
	"rpc-content/pkg/code"
	"context"
	"time"

	"rpc-content/client/user/user"
	"rpc-content/content"
	"rpc-content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminPendingListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminPendingListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminPendingListLogic {
	return &AdminPendingListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminPendingListLogic) AdminPendingList(in *content.AdminPendingListRequest) (resp *content.AdminPendingListResponse, err error) {
	resp = new(content.AdminPendingListResponse)
	resp.Data = new(content.SearchData)
	resp.Data.Items = make([]*content.SearchItem, 0)

	if in.PageSize == 0 {
		in.PageSize = 20
	}

	articles, err := l.svcCtx.ArticleModel.ArticlesPending(l.ctx, int(in.PageSize)+1, in.Cursor)
	if err != nil {
		resp.Code = int64(code.ServerErr.Code())
		resp.Msg = err.Error()
		return resp, nil
	}

	var isEnd bool
	if len(articles) > int(in.PageSize) {
		articles = articles[:in.PageSize]
		isEnd = false
	} else {
		isEnd = true
	}

	// 转换待审核文章实体为响应数据项
	items := make([]*content.SearchItem, 0, len(articles))
	for _, v := range articles {
		items = append(items, &content.SearchItem{
			ArticleId:   v.Id,
			Title:       v.Title,
			Description: v.Description,
			Cover:       v.Cover,
			AuthorId:    v.AuthorId,
			LikeNum:     v.LikeNum,
			CommentNum:  v.CommentNum,
			PublishTime: v.PublishTime.Format("2006-01-02 15:04:05"),
		})
	}

	if l.svcCtx.UserRPC != nil {
		l.populateAuthorNames(items)
	}

	var cursor int64
	if len(items) > 0 && !isEnd {
		last := items[len(items)-1]
		t, err := time.ParseInLocation("2006-01-02 15:04:05", last.PublishTime, time.Local)
		if err == nil {
			cursor = t.Unix()
		}
	}

	resp.Data.Items = items
	resp.Data.Cursor = cursor
	resp.Data.IsEnd = isEnd
	return resp, nil
}

func (l *AdminPendingListLogic) populateAuthorNames(items []*content.SearchItem) {
	// 批量查询并填充文章作者名称
	for _, v := range items {
		u, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{UserId: v.AuthorId})
		if err != nil {
			continue
		}
		v.AuthorName = u.Username
	}
}
