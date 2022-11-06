package grpc_service

import (
	"context"
	configService "distributedConfig/internal/delivery/proto"
	"distributedConfig/internal/entity"
	"distributedConfig/internal/usecase"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ConfigService struct {
	configUseCase usecase.ConfigUseCase
}

func NewConfigService(configUseCase usecase.ConfigUseCase) *ConfigService {
	return &ConfigService{configUseCase: configUseCase}
}

func (s *ConfigService) CreateConfig(ctx context.Context, r *configService.Config) (*configService.ConfigResponse, error) {
	config := &entity.Config{
		Name:    r.ServiceName,
		Data:    r.Data,
		Version: 1,
	}
	err := s.configUseCase.CreateConfig(config)
	if err != nil && err == usecase.ErrConfigAlreadyExists {
		return nil, status.Errorf(409, "Unable to create %s config: %s", r.ServiceName, err)
	} else if err != nil {
		return nil, status.Errorf(500, "Unable to create %s config: %s", r.ServiceName, err)
	}
	return &configService.ConfigResponse{
		Config: &configService.Config{
			ServiceName: config.Name,
			Data:        config.Data,
		},
		Version:   1,
		CreatedAt: timestamppb.New(config.CreatedAt),
	}, nil
}

func (s *ConfigService) GetConfig(ctx context.Context, r *configService.ConfigName) (*configService.ConfigResponse, error) {
	config, err := s.configUseCase.GetConfig(r.ServiceName)
	if err != nil && err == usecase.ErrConfigNotFound {
		return nil, status.Errorf(404, "Unable to get %s config: %s", r.ServiceName, err)
	} else if err != nil {
		return nil, status.Errorf(500, "Unable to get %s config: %s", r.ServiceName, err)
	}
	return &configService.ConfigResponse{
		Config: &configService.Config{
			ServiceName: config.Name,
			Data:        config.Data,
		},
		Version:   config.Version,
		CreatedAt: timestamppb.New(config.CreatedAt),
	}, nil
}

func (s *ConfigService) GetConfigByVersion(ctx context.Context, r *configService.ConfigNameAndVersion) (*configService.ConfigResponse, error) {
	config, err := s.configUseCase.GetConfigByVersion(r.ServiceName, r.Version)
	if err != nil && err == usecase.ErrConfigNotFound {
		return nil, status.Errorf(404, "Unable to get %s config with version %d : %s", r.ServiceName, r.Version, err)
	} else if err != nil {
		return nil, status.Errorf(500, "Unable to get %s config with version %d: %s", r.ServiceName, r.Version, err)
	}
	return &configService.ConfigResponse{
		Config: &configService.Config{
			ServiceName: config.Name,
			Data:        config.Data,
		},
		Version:   config.Version,
		CreatedAt: timestamppb.New(config.CreatedAt),
	}, nil
}

func (s *ConfigService) UpdateConfig(ctx context.Context, r *configService.Config) (*configService.ConfigResponse, error) {
	config := &entity.Config{
		Name: r.ServiceName,
		Data: r.Data,
	}
	err := s.configUseCase.UpdateConfig(config)
	if err != nil && err == usecase.ErrConfigNotFound {
		return nil, status.Errorf(404, "Unable to update %s config: %s", r.ServiceName, err)
	} else if err != nil {
		return nil, status.Errorf(500, "Unable to update %s config: %s", r.ServiceName, err)
	}
	return &configService.ConfigResponse{
		Config: &configService.Config{
			ServiceName: config.Name,
			Data:        config.Data,
		},
		Version:   config.Version,
		CreatedAt: timestamppb.New(config.CreatedAt),
	}, nil
}

func (s *ConfigService) DeleteConfig(ctx context.Context, r *configService.ConfigName) (*configService.DeleteResponse, error) {
	var err error
	err = s.configUseCase.DeleteConfig(r.ServiceName)
	if err != nil && err == usecase.ErrConfigWasRecentlyUsed {
		return nil, status.Errorf(403, "Unable to delete %s config: %s", r.ServiceName, err)
	} else if err != nil {
		return nil, status.Errorf(500, "Unable to delete %s config: %s", r.ServiceName, err)
	} else {
		return &configService.DeleteResponse{
			Message: "Config was deleted",
		}, nil
	}
}

func (s *ConfigService) DeleteConfigVersion(ctx context.Context, r *configService.ConfigNameAndVersion) (*configService.DeleteResponse, error) {
	var err error
	err = s.configUseCase.DeleteConfigVersion(r.ServiceName, r.Version)
	if err != nil && err == usecase.ErrConfigWasRecentlyUsed {
		return nil, status.Errorf(403, "Unable to delete %s config with version %d: %s", r.ServiceName, r.Version, err)
	} else if err != nil {
		return nil, status.Errorf(500, "Unable to delete %s config with version %d: %s", r.ServiceName, r.Version, err)
	} else {
		return &configService.DeleteResponse{
			Message: "Config was deleted",
		}, nil
	}
}

func (s *ConfigService) ListConfigs(r *configService.ListRequest, stream configService.ConfigService_ListConfigsServer) error {
	configs, err := s.configUseCase.GetConfigs(r.ServiceName)
	if err != nil && err == usecase.ErrConfigNotFound {
		return status.Errorf(404, "Unable to get %s configs: %s", r.ServiceName, err)
	} else if err != nil {
		return status.Errorf(500, "Unable to get %s configs: %s", r.ServiceName, err)
	}
	for _, config := range configs {
		err := stream.Send(&configService.ConfigResponse{
			Config: &configService.Config{
				ServiceName: config.Name,
				Data:        config.Data,
			},
			Version:   config.Version,
			CreatedAt: timestamppb.New(config.CreatedAt),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ConfigService) SetRelevantConfig(ctx context.Context, r *configService.ConfigNameAndVersion) (*configService.ConfigResponse, error) {
	config, err := s.configUseCase.SetRelevantConfig(r.ServiceName, r.Version)
	if err != nil && err == usecase.ErrConfigNotFound {
		return nil, status.Errorf(404, "Unable to set relevant %s config: %s", r.ServiceName, err)
	} else if err != nil {
		return nil, status.Errorf(500, "Unable to set relevant %s config: %s", r.ServiceName, err)
	}
	return &configService.ConfigResponse{
		Config: &configService.Config{
			ServiceName: config.Name,
			Data:        config.Data,
		},
		Version:   config.Version,
		CreatedAt: timestamppb.New(config.CreatedAt),
	}, nil
}
