package internal

import (
	"context"
	"distributedConfig/config"
	grpc_service "distributedConfig/internal/delivery/grpc"
	"distributedConfig/internal/delivery/proto"
	"distributedConfig/pkg/logger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net"
	"net/http"
)

func RunGatewayServer(configService *grpc_service.ConfigService, cfg *config.Config, l *logger.Logger) {
	grpcMux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := proto.RegisterConfigServiceHandlerServer(ctx, grpcMux, configService)
	if err != nil {
		l.Fatal("Failed to register gateway: %v", err)
		return
	}
	l.Info("Starting gRPC gateway on port %s", cfg.Server.GatewayPort)
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	listener, err := net.Listen("tcp", ":"+cfg.Server.GatewayPort)
	if err != nil {
		l.Fatal("Failed to listen: %v", err)
		return
	}
	if err := http.Serve(listener, mux); err != nil {
		l.Fatal("Failed to serve: %v", err)
		return
	}
}
