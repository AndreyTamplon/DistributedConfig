package interceptors

import (
	"context"
	"distributedConfig/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

type Interceptor struct {
	l logger.Logger
}

func NewInterceptor(l logger.Logger) *Interceptor {
	return &Interceptor{l: l}
}

func (i *Interceptor) Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx, req)
	i.l.Info("Method: %s, Time: %v, Metadata: %v, Err: %v", info.FullMethod, time.Since(start), md, err)
	return reply, err
}
