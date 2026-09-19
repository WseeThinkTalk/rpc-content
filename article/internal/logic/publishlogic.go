package logic

import (
	"rpc-content/article/internal/code"
	"rpc-content/article/internal/model"
	"rpc-content/article/internal/types"
	"rpc-content/tag/tag"
	"context"
	"strconv"
	"time"

	"rpc-content/article/internal/svc"
	"rpc-content/article/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

const (
	maxTitleLength       = 200
	maxContentLength     = 100000
	maxDescriptionLength = 500
)

func (l *PublishLogic) Publish(in *pb.PublishRequest) (*pb.PublishResponse, error) {
	if in.UserId <= 0 {
		return nil, code.UserIdInvalid
	}
	if len(in.Title) == 0 {
		return nil, code.ArticleTitleCantEmpty
	}
	if len(in.Title) > maxTitleLength {
		return nil, code.ArticleTitleTooLong
	}
	if len(in.Content) == 0 {
		return nil, code.ArticleContentCantEmpty
	}
	if len(in.Content) > maxContentLength {
		return nil, code.ArticleContentTooLong
	}
	if len(in.Description) > maxDescriptionLength {
		return nil, code.ArticleDescTooLong
	}
	ret, err := l.svcCtx.ArticleMOdel.Insert(l.ctx, &model.Article{
		AuthorId:    in.UserId,
		Title:       in.Title,
		Content:     in.Content,
		Description: in.Description,
		Cover:       in.Cover,
		Status:      types.ArticleStatusPending, // 修改为待审核状态
		PublishTime: time.Now(),
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	})

	if err != nil {
		l.Logger.Errorf("Publish Insert req: %v error: %v", in, err)
		return nil, err
	}

	articleId, err := ret.LastInsertId()
	if err != nil {
		l.Logger.Errorf("LastInsertId error: %v", err)
		return nil, err
	}

	var (
		articleIdStr   = strconv.FormatInt(articleId, 10)
		publishTimeKey = articlesKey(in.UserId, types.SortPublishTime)
		likeNumKey     = articlesKey(in.UserId, types.SortLikeCount)
	)
	b, _ := l.svcCtx.BizRedis.ExistsCtx(l.ctx, publishTimeKey)
	if b {
		_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, publishTimeKey, time.Now().Unix(), articleIdStr)
		if err != nil {
			logx.Errorf("ZaddCtx req: %v error: %v", in, err)
		}
	}
	b, _ = l.svcCtx.BizRedis.ExistsCtx(l.ctx, likeNumKey)
	if b {
		_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, likeNumKey, 0, articleIdStr)
		if err != nil {
			logx.Errorf("ZaddCtx req: %v error: %v", in, err)
		}
	}

	// 关联标签
	if len(in.TagIds) > 0 {
		for _, tagId := range in.TagIds {
			_, err := l.svcCtx.TagRPC.TagResource(l.ctx, &tag.TagResourceRequest{
				BizId:    "article",
				TargetId: articleId,
				TagId:    tagId,
			})
			if err != nil {
				l.Logger.Errorf("Publish TagResource articleId=%d tagId=%d error: %v", articleId, tagId, err)
			}
		}
	}

	return &pb.PublishResponse{ArticleId: articleId}, nil
}
