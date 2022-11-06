package database

import (
	"database/sql"
	"distributedConfig/config"
	"fmt"
	_ "github.com/lib/pq"
)

func NewDB(c *config.Config) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Dbname,
		c.Database.Password,
	)

	db, err := sql.Open(c.Database.Driver, dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
