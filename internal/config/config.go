package config

import (
	"lignis/internal/service/auth"
	"lignis/internal/service/minio"
	"lignis/internal/storage"
)

type Config struct {
	Mongo *storage.Config
	Minio *minio.Config
	Auth  *auth.Config
}
