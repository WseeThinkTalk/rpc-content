package articlelogic

import (
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
		l.Logger.Errorf("[AdminPendingList] ArticlesPending error: %v", err)
		return nil, err
	}

	var isEnd bool
	if len(articles) > int(in.PageSize) {
		articles = articles[:in.PageSize]
		isEnd = false
	} else {
		isEnd = true
	}

	items := make([]*content.SearchItem, 0, len(articles))
	for _, art := range articles {
		items = append(items, &content.SearchItem{
			ArticleId:   art.Id,
			Title:       art.Title,
			Description: art.Description,
			Cover:       art.Cover,
			AuthorId:    art.AuthorId,
			LikeNum:     art.LikeNum,
			CommentNum:  art.CommentNum,
			PublishTime: art.PublishTime.Format("2006-01-02 15:04:05"),
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
	for _, item := range items {
		u, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{UserId: item.AuthorId})
		if err != nil {
			l.Logger.Errorf("[AdminPendingList] FindById userId: %d error: %v", item.AuthorId, err)
			continue
		}
		item.AuthorName = u.Username
	}
}
