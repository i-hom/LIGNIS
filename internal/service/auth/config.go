package auth

type Config struct {
	SignKey string `envconfig:"JWT_SIGN_KEY"`
	TTL     int    `envconfig:"JWT_TTL"`
}
