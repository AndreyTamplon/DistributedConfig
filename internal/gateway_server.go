package internal

import (
	"context"
	grpc_service "distributedConfig/internal/delivery/grpc"
	"distributedConfig/internal/delivery/proto"
	"distributedConfig/pkg/logger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net"
	"net/http"
)

func RunGatewayServer(configService *grpc_service.ConfigService, l *logger.Logger) {
	grpcMux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := proto.RegisterConfigServiceHandlerServer(ctx, grpcMux, configService)
	if err != nil {
		l.Fatal("Failed to register gateway: %v", err)
		return
	}
	l.Info("Starting gRPC gateway on port %s", 8085)
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	listener, err := net.Listen("tcp", ":8085")
	if err != nil {
		l.Fatal("Failed to listen: %v", err)
		return
	}
	if err := http.Serve(listener, mux); err != nil {
		l.Fatal("Failed to serve: %v", err)
		return
	}
}
