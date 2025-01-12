package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type DBConfig struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:"postgres"`
	DBName   string `yaml:"dbname" env:"DB_NAME" env-default:"postgres"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers" env:"KAFKA_BROKERS" env-separator:"," env-default:"localhost:9092"`
	Topic   string   `yaml:"topic" env:"KAFKA_TOPIC" env-default:"user-segments"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:":8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type AppConfig struct {
	DB     DBConfig    `yaml:"database"`
	Kafka  KafkaConfig `yaml:"kafka"`
	Server HTTPServer  `yaml:"http_server"`
}

func LoadDBConfig(configPath string) (*AppConfig, error) {
	var cfg AppConfig
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, err
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
