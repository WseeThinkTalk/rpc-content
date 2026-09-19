package main

import (
	"flag"
	"fmt"

	articlerpc "rpc-content/article"
	"rpc-content/pkg/env"
	"rpc-content/pkg/interceptors"
	qarpc "rpc-content/qa"
	tagrpc "rpc-content/tag"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/content.yaml", "the config file")

type Config struct {
	zrpc.RpcServerConf
	DataSource string
	CacheRedis cache.CacheConf
	BizRedis   redis.RedisConf
	DB         struct {
		DataSource   string
		MaxOpenConns int `json:",default=10"`
		MaxIdleConns int `json:",default=100"`
		MaxLifetime  int `json:",default=3600"`
	}
	UserRPC zrpc.RpcClientConf
	TagRPC  zrpc.RpcClientConf
	Es      struct {
		Addresses []string
		Username  string
		Password  string
	}
}

func main() {
	flag.Parse()

	env.LoadEnv()

	var c Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	if c.DB.DataSource == "" {
		c.DB.DataSource = c.DataSource
	}

	articleConf := articlerpc.Config{
		RpcServerConf: c.RpcServerConf,
		DataSource:    c.DataSource,
		CacheRedis:    c.CacheRedis,
		BizRedis:      c.BizRedis,
		UserRPC:       c.UserRPC,
		TagRPC:        c.TagRPC,
	}

	qaConf := qarpc.Config{
		RpcServerConf: c.RpcServerConf,
		DB:            c.DB,
		BizRedis:      c.BizRedis,
		Es:            c.Es,
	}

	tagConf := tagrpc.Config{
		RpcServerConf: c.RpcServerConf,
		DB:            c.DB,
		BizRedis:      c.BizRedis,
	}

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		articlerpc.Register(grpcServer, articleConf)
		qarpc.Register(grpcServer, qaConf)
		tagrpc.Register(grpcServer, tagConf)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(interceptors.ServerErrorInterceptor())
	defer s.Stop()

	fmt.Printf("Starting unified content rpc server (article, qa, tag) at %s...\n", c.ListenOn)
	s.Start()
}
