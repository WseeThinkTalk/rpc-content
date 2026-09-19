package svc

import (
	"rpc-content/qa/internal/config"
	"rpc-content/qa/internal/model"
	"rpc-content/pkg/es"
	"rpc-content/pkg/orm"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"golang.org/x/sync/singleflight"
)

type ServiceContext struct {
	Config            config.Config
	DB                *orm.DB
	QuestionModel     *model.QuestionModel
	AnswerModel       *model.AnswerModel
	BizRedis          *redis.Redis
	SingleFlightGroup singleflight.Group
	Es                *es.Es
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
		Config:        c,
		DB:            db,
		QuestionModel: model.NewQuestionModel(db.DB),
		AnswerModel:   model.NewAnswerModel(db.DB),
		BizRedis:      rds,
		Es: es.MustNewEs(&es.Config{
			Addresses: c.Es.Addresses,
			Username:  c.Es.Username,
			Password:  c.Es.Password,
		}),
	}
}
