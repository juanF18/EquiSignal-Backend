package config

type Config struct {
	HttpPort string
}

func LoadConfig() *Config {
	return &Config{
		HttpPort: "8080",
	}
}
