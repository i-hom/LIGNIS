package minio

type Config struct {
	Username string `envconfig:"MINIO_ROOT_USER"`
	Password string `envconfig:"MINIO_ROOT_PASSWORD"`
	Endpoint string `envconfig:"MINIO_ENDPOINT"`
	Bucket   string `envconfig:"MINIO_BUCKET_NAME"`
}
