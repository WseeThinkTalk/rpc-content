package qarpc

import (
	"rpc-content/qa/internal/config"
	"rpc-content/qa/internal/server"
	"rpc-content/qa/internal/svc"
	"rpc-content/qa/pb"

	"google.golang.org/grpc"
)

type Config = config.Config

func Register(grpcServer *grpc.Server, c Config) {
	ctx := svc.NewServiceContext(c)
	pb.RegisterQAServer(grpcServer, server.NewQAServer(ctx))
}
