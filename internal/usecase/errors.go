package usecase

import "errors"

var (
	ErrConfigWasRecentlyUsed = errors.New("config was recently used")
	ErrConfigNotFound        = errors.New("config not found")
	ErrConfigAlreadyExists   = errors.New("config already exists")
)
