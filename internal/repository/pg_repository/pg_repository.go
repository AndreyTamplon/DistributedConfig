package pg_repository

import (
	"database/sql"
	"distributedConfig/internal/entity"
	"distributedConfig/internal/usecase"
	"time"
)

type ConfigRepository struct {
	db *sql.DB
}

func NewConfigRepository(db *sql.DB) *ConfigRepository {
	return &ConfigRepository{db: db}
}

func (r *ConfigRepository) CreateConfig(config *entity.Config) error {
	if err := config.Validate(); err != nil {
		return err
	}
	err := r.db.QueryRow("INSERT INTO configs (name, version) VALUES ($1, $2) RETURNING id, version, created_at",
		config.Name, config.Version).Scan(&config.ID, &config.Version, &config.CreatedAt)
	if err != nil {
		return err
	}
	err = r.insertData(config.ID, config.Data)
	if err != nil {
		return err
	}
	_, err = r.SetRelevantConfig(config.Name, config.Version)
	return err
}

func (r *ConfigRepository) GetConfig(name string) (*entity.Config, error) {
	var config entity.Config
	err := r.db.QueryRow("SELECT id, name, version, created_at FROM configs WHERE name = $1 AND relevant = TRUE",
		name).Scan(&config.ID, &config.Name, &config.Version, &config.CreatedAt)
	if err != nil {
		return nil, err
	}

	config.Data, err = r.GetDataByConfigID(config.ID)
	if err != nil {
		return nil, err
	}
	err = r.updateLastUsed(config.ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, usecase.ErrConfigNotFound
	}
	return &config, nil
}

func (r *ConfigRepository) GetConfigs(name string) ([]*entity.Config, error) {
	rows, err := r.db.Query("SELECT id, name, version, created_at FROM configs WHERE name = $1 ORDER BY version DESC", name)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, usecase.ErrConfigNotFound
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
	}(rows)

	var configs []*entity.Config
	for rows.Next() {
		var config entity.Config
		err = rows.Scan(&config.ID, &config.Name, &config.Version, &config.CreatedAt)
		if err != nil {
			return nil, err
		}
		config.Data, err = r.GetDataByConfigID(config.ID)
		if err != nil {
			return nil, err
		}
		err := r.updateLastUsed(config.ID)
		if err != nil {
			return nil, err
		}
		configs = append(configs, &config)
	}
	return configs, nil
}

func (r *ConfigRepository) GetConfigByVersion(name string, version int64) (*entity.Config, error) {
	var config entity.Config
	err := r.db.QueryRow("SELECT id, name, version, created_at FROM configs WHERE name = $1 AND version = $2",
		name, version).Scan(&config.ID, &config.Name, &config.Version, &config.CreatedAt)
	if err != nil {
		return nil, err
	}
	config.Data, err = r.GetDataByConfigID(config.ID)
	if err != nil {
		return nil, err
	}
	err = r.updateLastUsed(config.ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, usecase.ErrConfigNotFound
	}
	return &config, nil
}

func (r *ConfigRepository) DeleteConfig(name string) error {
	_, err := r.db.Exec("DELETE FROM pairs WHERE config_id IN (SELECT id FROM configs WHERE name = $1)", name)
	if err != nil {
		return err
	}
	_, err = r.db.Exec("DELETE FROM configs WHERE name = $1", name)
	return err
}

func (r *ConfigRepository) DeleteConfigVersion(name string, version int64) error {
	_, err := r.db.Exec("DELETE FROM pairs WHERE config_id IN (SELECT id FROM configs WHERE name = $1 AND version = $2)",
		name, version)
	if err != nil {
		return err
	}
	_, err = r.db.Exec("DELETE FROM configs WHERE name = $1 AND version = $2", name, version)
	return err
}

func (r *ConfigRepository) UpdateConfig(config *entity.Config) error {
	if err := config.Validate(); err != nil {
		return err
	}
	version, err := r.GetLastVersion(config.Name)
	if err != nil && err != sql.ErrNoRows {
		return err
	} else if err == sql.ErrNoRows {
		config.Version = 1
	} else {
		config.Version = version + 1
	}
	err = r.db.QueryRow("INSERT INTO configs (name, version) VALUES ($1, $2) RETURNING id, version, created_at",
		config.Name, config.Version).Scan(&config.ID, &config.Version, &config.CreatedAt)
	if err != nil {
		return err
	}
	err = r.insertData(config.ID, config.Data)
	if err != nil {
		return err
	}
	_, err = r.SetRelevantConfig(config.Name, config.Version)
	return err
}

func (r *ConfigRepository) SetRelevantConfig(name string, version int64) (*entity.Config, error) {
	_, err := r.db.Exec("UPDATE configs SET relevant = FALSE WHERE name = $1", name)
	if err != nil {
		return nil, err
	}
	_, err = r.db.Exec("UPDATE configs SET relevant = TRUE, last_used = $1 WHERE name = $2 AND version = $3",
		time.Now(), name, version)
	if err != nil {
		return nil, err
	}
	return r.GetConfigByVersion(name, version)
}

func (r *ConfigRepository) GetDataByConfigID(id int) (map[string]string, error) {
	rows, err := r.db.Query("SELECT key, value FROM pairs WHERE config_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
	}(rows)
	data := make(map[string]string)
	for rows.Next() {
		var key, value string
		err = rows.Scan(&key, &value)
		if err != nil {
			return nil, err
		}
		data[key] = value
	}
	return data, nil
}

func (r *ConfigRepository) GetRelevantLastUsed(name string) (time.Time, error) {
	var lastUsed time.Time
	err := r.db.QueryRow("SELECT last_used FROM configs WHERE name = $1 AND relevant = TRUE", name).Scan(&lastUsed)
	if err == sql.ErrNoRows {
		return time.Time{}, usecase.ErrConfigNotFound
	}
	return lastUsed, err
}

func (r *ConfigRepository) GetLastUsedByVersion(name string, version int64) (time.Time, error) {
	var lastUsed time.Time
	err := r.db.QueryRow("SELECT last_used FROM configs WHERE name = $1 AND version = $2", name, version).Scan(&lastUsed)
	if err == sql.ErrNoRows {
		return time.Time{}, usecase.ErrConfigNotFound
	}
	return lastUsed, err
}

func (r *ConfigRepository) insertData(configID int, data map[string]string) error {
	for key, value := range data {
		_, err := r.db.Exec("INSERT INTO pairs (config_id, key, value) VALUES ($1, $2, $3)", configID, key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ConfigRepository) IsConfigExists(name string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS (SELECT 1 FROM configs WHERE name = $1)", name).Scan(&exists)
	return exists, err
}

func (r *ConfigRepository) IsConfigVersionExists(name string, version int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS (SELECT 1 FROM configs WHERE name = $1 AND version = $2)", name, version).Scan(&exists)
	return exists, err
}

func (r *ConfigRepository) IsConfigRelevant(name string, version int64) (bool, error) {
	var relevant bool
	err := r.db.QueryRow("SELECT relevant FROM configs WHERE name = $1 AND version = $2", name, version).Scan(&relevant)
	return relevant, err
}

func (r *ConfigRepository) GetLastVersion(name string) (int64, error) {
	var version int64
	err := r.db.QueryRow("SELECT version FROM configs WHERE name = $1 ORDER BY version DESC LIMIT 1", name).Scan(&version)
	return version, err
}

func (r *ConfigRepository) updateLastUsed(configID int) error {
	_, err := r.db.Exec("UPDATE configs SET last_used = $1 WHERE id = $2", time.Now(), configID)
	return err
}
