package entity

import "testing"

func TestConfig(t *testing.T) *Config {
	t.Helper()

	return &Config{
		Name: "test",
		Data: map[string]string{"key1": "value1", "key2": "value2", "key3": "value3"},
	}
}
