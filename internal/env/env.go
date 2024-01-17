package env

type Config struct {
	Port int `env:"PORT,default=50000"`
}
