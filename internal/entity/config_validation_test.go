package entity

import "testing"

func TestConfig_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		c       func() *Config
		isValid bool
	}{
		{
			name: "valid",
			c: func() *Config {
				return TestConfig(t)
			},
			isValid: true,
		},
		{
			name: "empty name",
			c: func() *Config {
				c := TestConfig(t)
				c.Name = ""

				return c
			},
			isValid: false,
		},
		{
			name: "empty data",
			c: func() *Config {
				c := TestConfig(t)
				c.Data = nil

				return c
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.c()

			err := c.Validate()

			if tc.isValid && err != nil {
				t.Errorf("expected config to be valid, got %v", err)
			}

			if !tc.isValid && err == nil {
				t.Errorf("expected config to be invalid")
			}
		})
	}
}
