package entity

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"time"
)

type Config struct {
	ID        int               `json:"id"`
	Name      string            `json:"name"`
	Data      map[string]string `json:"data"`
	CreatedAt time.Time         `json:"created_at"`
	Version   int64             `json:"version"`
}

func (config *Config) Validate() error {
	return validation.ValidateStruct(
		config,
		validation.Field(&config.Name, validation.Required, validation.Length(1, 255)),
		validation.Field(&config.Data, validation.Required),
	)
}
