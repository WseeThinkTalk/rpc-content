package svc

import (
	"rpc-content/client/user/user"
	"rpc-content/internal/config"
	articlemodel "rpc-content/internal/model/article"
	qamodel "rpc-content/internal/model/qa"
	tagmodel "rpc-content/internal/model/tag"
	"rpc-content/pkg/es"
	"rpc-content/pkg/orm"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"golang.org/x/sync/singleflight"
)

type ServiceContext struct {
	Config            config.Config
	DB                *orm.DB
	BizRedis          *redis.Redis
	SingleFlightGroup singleflight.Group
	Es                *es.Es
	UserRPC           user.User

	ArticleModel     articlemodel.ArticleModel
	QuestionModel    *qamodel.QuestionModel
	AnswerModel      *qamodel.AnswerModel
	TagModel         *tagmodel.TagModel
	TagResourceModel *tagmodel.TagResourceModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	dsn := c.DB.DataSource
	if dsn == "" {
		dsn = c.DataSource
	}

	db := orm.MustNewPostgres(&orm.Config{
		DSN:          dsn,
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

	artModel := articlemodel.NewArticleModel(db.DB)

	sc := &ServiceContext{
		Config:           c,
		DB:               db,
		BizRedis:         rds,
		ArticleModel:     artModel,
		QuestionModel:    qamodel.NewQuestionModel(db.DB),
		AnswerModel:      qamodel.NewAnswerModel(db.DB),
		TagModel:         tagmodel.NewTagModel(db.DB),
		TagResourceModel: tagmodel.NewTagResourceModel(db.DB),
	}

	if len(c.Es.Addresses) > 0 {
		sc.Es = es.MustNewEs(&es.Config{
			Addresses: c.Es.Addresses,
			Username:  c.Es.Username,
			Password:  c.Es.Password,
		})
	}
	if len(c.UserRPC.Endpoints) > 0 || c.UserRPC.Target != "" {
		sc.UserRPC = user.NewUser(zrpc.MustNewClient(c.UserRPC))
	}

	return sc
}
