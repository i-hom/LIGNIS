package storage

type Config struct {
	MongoURI      string `envconfig:"MONGO_URI"`
	MongoDatabase string `envconfig:"MONGO_DB"`
}
