package app

import (
	"distributedConfig/config"
	"distributedConfig/internal"
	"distributedConfig/internal/delivery/grpc"
	"distributedConfig/internal/repository/pg_repository"
	"distributedConfig/internal/usecase"
	"distributedConfig/pkg/database"
	"distributedConfig/pkg/logger"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Logger.LogLevel)
	db, err := database.NewDB(cfg)
	if err != nil {
		l.Fatal("Failed to connect to database: %v", err)
		return
	}
	defer db.Close()
	l.Info("Database connected")
	configRepository := pg_repository.NewConfigRepository(db)
	configUseCase := usecase.NewConfigUseCase(*l, configRepository, cfg)
	configService := grpc_service.NewConfigService(*configUseCase)
	go internal.RunGatewayServer(configService, cfg, l)
	internal.RunGrpcServer(configService, cfg, l)

}
