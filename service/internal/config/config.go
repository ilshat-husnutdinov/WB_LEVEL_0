package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DB         DB         `yaml:"database"`
	HTTPServer HTTPServer `yaml:"http_server"`
	STAN       STAN       `yaml:"stan"`
}

type DB struct {
	DSN string `yaml:"DSN"`
}

func (DB *DB) GetDSN() string {
	return DB.DSN
}

type HTTPServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type STAN struct {
	ClusterID string `yaml:"cluster_id"`
	ClientID  string `yaml:"client_id"`
	Subject   string `yaml:"subject"`
	Durable   string `yaml:"durable"`
}

// LoadConfigYaml читает конфигурацию по указанному пути и возвращает *Config
func LoadConfigYaml(configYaml string) *Config {
	const op = "internal.config.ReadConfigYaml"

	var cfg Config

	yamlFile, err := os.ReadFile(configYaml)
	if err != nil {
		log.Fatalf("%s: can`t read a file. Reason: %v", op, err)
	}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("%s: can`t unmarshal a yaml-file. Reason: %v", op, err)
	}

	// проверка на корректное чтение конфига
	if (cfg == Config{}) {
		log.Fatalf("%s: data from config wasn`t read", op)
	}
	return &cfg
}
