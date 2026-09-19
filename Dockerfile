FROM golang:1.26-alpine AS builder

WORKDIR /build

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/rpc-content .

FROM alpine:3.20

ENV TZ=Asia/Shanghai
RUN apk add --no-cache tzdata ca-certificates

WORKDIR /app
COPY --from=builder /app/rpc-content /app/rpc-content
COPY etc /app/etc

EXPOSE 9090

ENTRYPOINT ["/app/rpc-content"]
CMD ["-f", "etc/content.yaml"]
