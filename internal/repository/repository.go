package repository

import (
	"distributedConfig/internal/entity"
	"time"
)

type ConfigRepository interface {
	CreateConfig(config *entity.Config) error
	GetConfig(name string) (*entity.Config, error)
	GetConfigs(name string) ([]*entity.Config, error)
	GetConfigByVersion(name string, version int64) (*entity.Config, error)
	DeleteConfig(name string) error
	DeleteConfigVersion(name string, version int64) error
	UpdateConfig(config *entity.Config) error
	SetRelevantConfig(name string, version int64) (*entity.Config, error)
	GetRelevantLastUsed(name string) (time.Time, error)
	GetLastUsedByVersion(name string, version int64) (time.Time, error)
	IsConfigExists(name string) (bool, error)
	IsConfigVersionExists(name string, version int64) (bool, error)
	IsConfigRelevant(name string, version int64) (bool, error)
	GetLastVersion(name string) (int64, error)
}
