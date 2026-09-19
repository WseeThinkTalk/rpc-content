package svc

import (
	"rpc-content/article/internal/config"
	"rpc-content/article/internal/model"
	"rpc-content/tag/tag"
	"rpc-content/client/user/user"
	"rpc-content/pkg/orm"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"golang.org/x/sync/singleflight"
)

type ServiceContext struct {
	Config            config.Config
	ArticleMOdel      model.ArticleModel
	DB                *orm.DB
	BizRedis          *redis.Redis
	SingleFlightGroup singleflight.Group
	UserRPC           user.User
	TagRPC            tag.Tag
}

func NewServiceContext(c config.Config) *ServiceContext {
	rds, err := redis.NewRedis(c.BizRedis)
	if err != nil {
		panic(err)
	}
	db := orm.MustNewPostgres(&orm.Config{
		DSN:          c.DataSource,
		MaxOpenConns: 10,
		MaxIdleConns: 100,
		MaxLifetime:  3600,
	})
	return &ServiceContext{
		Config:       c,
		ArticleMOdel: model.NewArticleModel(db.DB),
		DB:           db,
		BizRedis:     rds,
		UserRPC:      user.NewUser(zrpc.MustNewClient(c.UserRPC)),
		TagRPC:       tag.NewTag(zrpc.MustNewClient(c.TagRPC)),
	}
}
