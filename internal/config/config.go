package config

import (
	"lignis/internal/service/auth"
	"lignis/internal/storage"
)

type Config struct {
	Mongo *storage.Config
	Auth  *auth.Config
}
