package configs

type Config struct {
	Servidor Servidor
	Balanca  *Balanca
}

type Servidor struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	ReadTimeout  int    `json:"readTimeout"`
	WriteTimeout int    `json:"writeTimeout"`
}

type Balanca struct {
	Modelo    string `json:"modelo"`
	Protocolo string `json:"protocolo"`
	Baud      int    `json:"baud"`
	Paridade  int    `json:"paridade"`
}

func NovoConfig() *Config {
	cfg := &Config{Balanca: &Balanca{}}
	cfg.Servidor.Host = "127.0.0.1"
	cfg.Servidor.Port = 8082
	cfg.Servidor.WriteTimeout = 10000
	cfg.Servidor.ReadTimeout = 5000
	return cfg
}
