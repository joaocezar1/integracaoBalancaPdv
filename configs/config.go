package configs

var cfg *Config

type Config struct {
	Host         string
	Port         int
	ReadTimeout  int
	WriteTimeout int
}

func NovoConfig() *Config {
	cfg := &Config{}
	cfg.Host = "127.0.0.1"
	cfg.Port = 8082
	cfg.WriteTimeout = 1000
	cfg.ReadTimeout = 1000
	return cfg
}
