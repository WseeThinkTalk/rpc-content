package svc

import (
	"rpc-content/tag/internal/config"
	"rpc-content/tag/internal/model"
	"rpc-content/pkg/orm"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config           config.Config
	DB               *orm.DB
	TagModel         *model.TagModel
	TagResourceModel *model.TagResourceModel
	BizRedis         *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := orm.MustNewMysql(&orm.Config{
		DSN:          c.DB.DataSource,
		MaxOpenConns: c.DB.MaxOpenConns,
		MaxIdleConns: c.DB.MaxIdleConns,
		MaxLifetime:  c.DB.MaxLifetime,
	})

	rds := redis.MustNewRedis(redis.RedisConf{
		Host:        c.BizRedis.Host,
		Pass:        c.BizRedis.Pass,
		Type:        c.BizRedis.Type,
		PingTimeout: 10000000000,
	})

	return &ServiceContext{
		Config:           c,
		DB:               db,
		TagModel:         model.NewTagModel(db.DB),
		TagResourceModel: model.NewTagResourceModel(db.DB),
		BizRedis:         rds,
	}
}
