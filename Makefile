## rpc 生成命令
rpc:
	goctl rpc protoc content.proto --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go_out=./content --go-grpc_out=./content --zrpc_out=./ -m --verbose --style=go_zero

run:
	go run content.go

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/rpc-content content.go
