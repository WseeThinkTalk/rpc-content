package tagrpc

import (
	"rpc-content/tag/internal/config"
	"rpc-content/tag/internal/server"
	"rpc-content/tag/internal/svc"
	"rpc-content/tag/pb"

	"google.golang.org/grpc"
)

type Config = config.Config

func Register(grpcServer *grpc.Server, c Config) {
	ctx := svc.NewServiceContext(c)
	pb.RegisterTagServer(grpcServer, server.NewTagServer(ctx))
}
