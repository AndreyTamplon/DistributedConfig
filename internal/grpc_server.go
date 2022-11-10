package internal

import (
	"distributedConfig/config"
	"distributedConfig/internal/delivery/grpc"
	"distributedConfig/internal/delivery/proto"
	"distributedConfig/internal/interceptors"
	"distributedConfig/pkg/logger"
	"google.golang.org/grpc"
	"net"
)

func RunGrpcServer(configService *grpc_service.ConfigService, cfg *config.Config, l *logger.Logger) {
	interceptor := interceptors.NewInterceptor(*l)
	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Logger))
	proto.RegisterConfigServiceServer(server, configService)
	l.Info("Starting gRPC server on port %s", cfg.Server.GPRCPort)
	listener, err := net.Listen("tcp", ":"+cfg.Server.GPRCPort)
	if err != nil {
		l.Fatal("Failed to listen: %v", err)
		return
	}
	if err := server.Serve(listener); err != nil {
		l.Fatal("Failed to serve: %v", err)
		return
	}
}
