package usecase

import (
	cfg "distributedConfig/config"
	"distributedConfig/internal/entity"
	"distributedConfig/internal/repository"
	"distributedConfig/pkg/logger"
	"time"
)

type ConfigUseCase struct {
	l          logger.Logger
	repository repository.ConfigRepository
	cfg        *cfg.Config
}

func NewConfigUseCase(l logger.Logger, repository repository.ConfigRepository, cfg *cfg.Config) *ConfigUseCase {
	return &ConfigUseCase{l: l, repository: repository, cfg: cfg}
}

func (c *ConfigUseCase) CreateConfig(config *entity.Config) error {
	exists, err := c.repository.IsConfigExists(config.Name)
	if err != nil {
		c.l.Error("Unable to check if config %s exists: %s", config.Name, err)
		return err
	}
	if exists {
		c.l.Error("Config %s already exists", config.Name)
		return ErrConfigAlreadyExists
	}
	err = c.repository.CreateConfig(config)
	if err != nil {
		c.l.Error("Unable to create config: %s", err)
		return err
	}
	c.l.Info("Config created: %s %d", config.Name, config.Version)
	return nil
}

func (c *ConfigUseCase) GetConfig(name string) (*entity.Config, error) {
	config, err := c.repository.GetConfig(name)
	if err != nil {
		c.l.Error("Unable to get config: %s", err)
		return nil, err
	}
	c.l.Info("Config got: %s %d", config.Name, config.Version)
	return config, nil
}

func (c *ConfigUseCase) GetConfigs(name string) ([]*entity.Config, error) {
	configs, err := c.repository.GetConfigs(name)
	if err != nil {
		c.l.Error("Unable to get configs: %s", err)
		return nil, err
	}
	c.l.Info("Configs got: %s", name)
	return configs, nil
}

func (c *ConfigUseCase) GetConfigByVersion(name string, version int64) (*entity.Config, error) {
	config, err := c.repository.GetConfigByVersion(name, version)
	if err != nil {
		c.l.Error("Unable to get config: %s with version %d", name, version)
		return nil, err
	}
	c.l.Info("Config got: %s %d", config.Name, config.Version)
	return config, nil
}

func (c *ConfigUseCase) DeleteConfig(name string) error {
	lastUsed, err := c.repository.GetRelevantLastUsed(name)
	if err != nil {
		c.l.Error("Unable to get time of last use of %s config: %s", name, err)
		return err
	}
	if !c.cfg.Server.DeleteConfigIfRecentlyUsed &&
		time.Now().Sub(lastUsed) < time.Duration(c.cfg.Server.RecentUseDurationDays)*24*time.Hour {
		c.l.Error("Unable to delete %s config: last use was less than 5 days ago", name)
		return ErrConfigWasRecentlyUsed
	} else {
		exists, err := c.repository.IsConfigExists(name)
		if err != nil {
			c.l.Error("Unable to check if config %s exists: %s", name, err)
			return err
		}
		if !exists {
			c.l.Error("Config %s not found", name)
			return ErrConfigNotFound
		}
		err = c.repository.DeleteConfig(name)
		if err != nil {
			c.l.Error("Unable to delete %s config: %s", name, err)
			return err
		}
		c.l.Info("Config deleted: %s", name)
		return nil
	}
}

func (c *ConfigUseCase) DeleteConfigVersion(name string, version int64) error {
	lastUsed, err := c.repository.GetLastUsedByVersion(name, version)
	if err != nil {
		c.l.Error("Unable to get time of last use of %s config: %s", name, err)
		return err
	}
	if !c.cfg.Server.DeleteConfigIfRecentlyUsed &&
		time.Now().Sub(lastUsed) < time.Duration(c.cfg.Server.RecentUseDurationDays)*24*time.Hour {
		c.l.Error("Unable to delete %s config: last use was less than 5 days ago", name)
		return ErrConfigWasRecentlyUsed
	} else {
		exists, err := c.repository.IsConfigVersionExists(name, version)
		if err != nil {
			c.l.Error("Unable to check if config %s with version %d exists: %s", name, version, err)
			return err
		}
		if !exists {
			c.l.Error("Config %s with version %d not found", name, version)
			return ErrConfigNotFound
		}
		isRelevant, err := c.repository.IsConfigRelevant(name, version)
		if err != nil {
			c.l.Error("Unable to check if config %s with version %d is relevant: %s", name, version, err)
			return err
		}
		if isRelevant {
			err := c.setNewRelevantBeforeDeletion(name)
			if err != nil {
				return err
			}
		}
		err = c.repository.DeleteConfigVersion(name, version)
		if err != nil {
			c.l.Error("Unable to delete %d version of %s config: %s", version, name, err)
			return err
		}
		c.l.Info("Version %d of config %s deleted", version, name)
		return nil
	}
}

func (c *ConfigUseCase) UpdateConfig(config *entity.Config) error {
	exists, err := c.repository.IsConfigExists(config.Name)
	if err != nil {
		c.l.Error("Unable to check if config %s exists: %s", config.Name, err)
		return err
	} else if !exists {
		c.l.Error("Config %s not found", config.Name)
		return ErrConfigNotFound
	}
	err = c.repository.UpdateConfig(config)
	if err != nil {
		c.l.Error("Unable to update config: %s", err)
		return err
	}
	c.l.Info("Config updated: %s %d", config.Name, config.Version)
	return nil
}

func (c *ConfigUseCase) setNewRelevantBeforeDeletion(name string) error {
	version, err := c.repository.GetLastVersion(name)
	if err != nil {
		c.l.Error("Unable to get last version of %s config: %s", name, err)
		return err
	}
	_, err = c.SetRelevantConfig(name, version)
	if err != nil {
		return err
	}
	return nil
}

func (c *ConfigUseCase) SetRelevantConfig(name string, version int64) (*entity.Config, error) {
	exists, err := c.repository.IsConfigVersionExists(name, version)
	if err != nil {
		c.l.Error("Unable to check if config %s with version %d exists: %s", name, version, err)
		return nil, err
	}
	if !exists {
		c.l.Error("Config %s with version %d not found", name, version)
		return nil, ErrConfigNotFound
	}
	config, err := c.repository.SetRelevantConfig(name, version)
	if err != nil {
		c.l.Error("Unable to get config: %s with version %d", name, version)
		return nil, err
	}
	c.l.Info("Version %d of config %s set relevant", version, name)

	return config, nil
}
