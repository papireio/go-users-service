package env

type Config struct {
	Port     int    `env:"PORT,default=50000"`
	MongoURL string `env:"MONGO_URL,default=mongodb://localhost:27017"`
}
