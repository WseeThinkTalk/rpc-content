package articlerpc

import (
	"rpc-content/article/internal/config"
	articleServer "rpc-content/article/internal/server"
	"rpc-content/article/internal/svc"
	"rpc-content/article/pb"

	"google.golang.org/grpc"
)

type Config = config.Config

func Register(grpcServer *grpc.Server, c Config) {
	ctx := svc.NewServiceContext(c)
	pb.RegisterArticleServer(grpcServer, articleServer.NewArticleServer(ctx))
}
