package pg_repository

import (
	"database/sql/driver"
	"distributedConfig/internal/entity"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestConfigRepository_CreateConfig(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()
	mock.ExpectQuery("INSERT INTO configs (name, version) VALUES ($1, $2) RETURNING id, version, created_at").
		WithArgs("test", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "version", "created_at"}).
			AddRow(1, 1, time.Now()))
	mock.ExpectExec("INSERT INTO pairs (config_id, key, value) VALUES ($1, $2, $3)").
		WithArgs(1, "key1", "value1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE configs SET relevant = FALSE WHERE name = $1").
		WithArgs("test").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE configs SET relevant = TRUE, last_used = $1 WHERE name = $2 AND version = $3").
		WithArgs(AnyTime{}, "test", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	rows := sqlmock.NewRows([]string{"id", "name", "version", "created_at"}).
		AddRow(1, "test", 1, time.Now())
	mock.ExpectQuery("SELECT id, name, version, created_at FROM configs WHERE name = $1 AND version = $2").
		WithArgs("test", 1).
		WillReturnRows(rows)
	pairRows := sqlmock.NewRows([]string{"key", "value"}).
		AddRow("key1", "value1").
		AddRow("key2", "value2")
	mock.ExpectQuery("SELECT key, value FROM pairs WHERE config_id = $1").
		WithArgs(1).
		WillReturnRows(pairRows)
	mock.ExpectExec("UPDATE configs SET last_used = $1 WHERE id = $2").
		WithArgs(AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	repo := NewConfigRepository(db)
	err = repo.CreateConfig(&entity.Config{
		Name:    "test",
		Version: 1,
		Data:    map[string]string{"key1": "value1"},
	})
	require.NoError(t, err)
}

func TestConfigRepository_GetConfig(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "version", "created_at"}).
		AddRow(1, "test", 1, time.Now())
	mock.ExpectQuery("SELECT id, name, version, created_at FROM configs WHERE name = $1 AND relevant = TRUE").
		WithArgs("test").
		WillReturnRows(rows)
	pairRows := sqlmock.NewRows([]string{"key", "value"}).
		AddRow("key1", "value1").
		AddRow("key2", "value2")
	mock.ExpectQuery("SELECT key, value FROM pairs WHERE config_id = $1").
		WithArgs(1).
		WillReturnRows(pairRows)
	mock.ExpectExec("UPDATE configs SET last_used = $1 WHERE id = $2").
		WithArgs(AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewConfigRepository(db)
	config, err := repo.GetConfig("test")
	require.NotNil(t, config)
	require.Equal(t, "test", config.Name)
	require.NoError(t, err)
}

func TestConfigRepository_GetConfigs(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "version", "created_at"}).
		AddRow(1, "test", 1, time.Now()).
		AddRow(2, "test", 2, time.Now())
	mock.ExpectQuery("SELECT id, name, version, created_at FROM configs WHERE name = $1 ORDER BY version DESC").
		WithArgs("test").
		WillReturnRows(rows)
	pairRows := sqlmock.NewRows([]string{"key", "value"}).
		AddRow("key1", "value1").
		AddRow("key2", "value2")
	mock.ExpectQuery("SELECT key, value FROM pairs WHERE config_id = $1").
		WithArgs(1).
		WillReturnRows(pairRows)
	mock.ExpectExec("UPDATE configs SET last_used = $1 WHERE id = $2").
		WithArgs(AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT key, value FROM pairs WHERE config_id = $1").
		WithArgs(2).
		WillReturnRows(pairRows)
	mock.ExpectExec("UPDATE configs SET last_used = $1 WHERE id = $2").
		WithArgs(AnyTime{}, 2).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewConfigRepository(db)
	configs, err := repo.GetConfigs("test")
	require.NotNil(t, configs)
	require.Equal(t, 2, len(configs))
	require.NoError(t, err)
}

func TestConfigRepository_GetConfigByVersion(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "version", "created_at"}).
		AddRow(1, "test", 1, time.Now())
	mock.ExpectQuery("SELECT id, name, version, created_at FROM configs WHERE name = $1 AND version = $2").
		WithArgs("test", 1).
		WillReturnRows(rows)
	pairRows := sqlmock.NewRows([]string{"key", "value"}).
		AddRow("key1", "value1").
		AddRow("key2", "value2")
	mock.ExpectQuery("SELECT key, value FROM pairs WHERE config_id = $1").
		WithArgs(1).
		WillReturnRows(pairRows)
	mock.ExpectExec("UPDATE configs SET last_used = $1 WHERE id = $2").
		WithArgs(AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewConfigRepository(db)
	config, err := repo.GetConfigByVersion("test", 1)
	require.NotNil(t, config)
	require.Equal(t, "test", config.Name)
	require.Equal(t, int64(1), config.Version)
	require.NoError(t, err)
}

func TestConfigRepository_DeleteConfig(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()
	mock.ExpectExec("DELETE FROM pairs WHERE config_id IN (SELECT id FROM configs WHERE name = $1)").
		WithArgs("test").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM configs WHERE name = $1").
		WithArgs("test").
		WillReturnResult(sqlmock.NewResult(1, 1))
	repo := NewConfigRepository(db)
	err = repo.DeleteConfig("test")
	require.NoError(t, err)
}
func TestConfigRepository_DeleteConfigByVersion(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()
	mock.ExpectExec("DELETE FROM pairs WHERE config_id IN (SELECT id FROM configs WHERE name = $1 AND version = $2)").
		WithArgs("test", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM configs WHERE name = $1 AND version = $2").
		WithArgs("test", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	repo := NewConfigRepository(db)
	err = repo.DeleteConfigVersion("test", 1)
	require.NoError(t, err)
}

func TestConfigRepository_GetLastVersion(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()
	rows := sqlmock.NewRows([]string{"version"}).
		AddRow(1)
	mock.ExpectQuery("SELECT version FROM configs WHERE name = $1 ORDER BY version DESC LIMIT 1").
		WithArgs("test").
		WillReturnRows(rows)
	repo := NewConfigRepository(db)
	version, err := repo.GetLastVersion("test")
	require.Equal(t, int64(1), version)
	require.NoError(t, err)
}

func TestConfigRepository_GetDataByConfigID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()
	rows := sqlmock.NewRows([]string{"key", "value"}).
		AddRow("key1", "value1").
		AddRow("key2", "value2")
	mock.ExpectQuery("SELECT key, value FROM pairs WHERE config_id = $1").
		WithArgs(1).
		WillReturnRows(rows)
	repo := NewConfigRepository(db)
	data, err := repo.GetDataByConfigID(1)
	require.Equal(t, 2, len(data))
	require.NoError(t, err)
}

func TestConfigRepository_GetLastUsedByVersion(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()
	rows := sqlmock.NewRows([]string{"last_used"}).
		AddRow(time.Now())
	mock.ExpectQuery("SELECT last_used FROM configs WHERE name = $1 AND version = $2").
		WithArgs("test", 1).
		WillReturnRows(rows)
	repo := NewConfigRepository(db)
	lastUsed, err := repo.GetLastUsedByVersion("test", 1)
	require.NotNil(t, lastUsed)
	require.NoError(t, err)
}

func TestConfigRepository_GetRelevantLastUsed(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()
	rows := sqlmock.NewRows([]string{"last_used"}).
		AddRow(time.Now())
	mock.ExpectQuery("SELECT last_used FROM configs WHERE name = $1 AND relevant = TRUE").
		WithArgs("test").
		WillReturnRows(rows)
	repo := NewConfigRepository(db)
	lastUsed, err := repo.GetRelevantLastUsed("test")
	require.NotNil(t, lastUsed)
	require.NoError(t, err)
}

func TestConfigRepository_IsConfigExists(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()
	mock.ExpectQuery("SELECT EXISTS (SELECT 1 FROM configs WHERE name = $1)").
		WithArgs("test").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	repo := NewConfigRepository(db)
	exists, err := repo.IsConfigExists("test")
	require.True(t, exists)
	require.NoError(t, err)
}

func TestConfigRepository_IsConfigExistsByVersion(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()
	mock.ExpectQuery("SELECT EXISTS (SELECT 1 FROM configs WHERE name = $1 AND version = $2)").
		WithArgs("test", 1).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	repo := NewConfigRepository(db)
	exists, err := repo.IsConfigVersionExists("test", 1)
	require.True(t, exists)
	require.NoError(t, err)
}

func TestConfigRepository_IsConfigRelevant(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()
	mock.ExpectQuery("SELECT relevant FROM configs WHERE name = $1 AND version = $2").
		WithArgs("test", 1).
		WillReturnRows(sqlmock.NewRows([]string{"relevant"}).AddRow(true))
	repo := NewConfigRepository(db)
	relevant, err := repo.IsConfigRelevant("test", 1)
	require.True(t, relevant)
	require.NoError(t, err)
}
