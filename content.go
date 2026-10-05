package main

import (
	"context"
	"flag"
	"fmt"

	"rpc-content/content"
	"rpc-content/internal/config"
	articleserver "rpc-content/internal/server/article"
	qaserver "rpc-content/internal/server/qa"
	tagserver "rpc-content/internal/server/tag"
	"rpc-content/internal/svc"
	"rpc-content/pkg/env"
	"rpc-content/pkg/lib/zapx"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/content.yaml", "the config file")

func main() {
	flag.Parse()

	env.LoadEnv()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	if c.DB.DataSource == "" {
		c.DB.DataSource = c.DataSource
	}

	// init logger
	writer, err := zapx.NewZapWriter()
	if err == nil {
		logx.SetWriter(writer)
	}

	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		registerServer(ctx, grpcServer)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(unaryServerInterceptor())
	defer s.Stop()

	fmt.Printf("Starting unified content rpc server (article, qa, tag) at %s...\n", c.ListenOn)
	s.Start()
}

// registerServer 注册 RPC 服务
func registerServer(ctx *svc.ServiceContext, grpcServer grpc.ServiceRegistrar) {
	content.RegisterArticleServer(grpcServer, articleserver.NewArticleServer(ctx))
	content.RegisterQAServer(grpcServer, qaserver.NewQAServer(ctx))
	content.RegisterTagServer(grpcServer, tagserver.NewTagServer(ctx))
}

// unaryServerInterceptor grpc 拦截器
func unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logx.Errorf("[RPC Panic] info: %s, recover: %v", info.FullMethod, r)
			}
		}()

		resp, err := handler(ctx, req)
		return resp, err
	}
}

