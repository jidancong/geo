package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Kafka         `yaml:"kafka"`
		Elasticsearch `yaml:"es"`
		ChunZhen      `yaml:"chunzhen"`
		IP2location   `yaml:"ip2location"`
	}

	Kafka struct {
		Host          string `yaml:"host"   env:"KAFKA_HOST"`
		Topic         string `yaml:"topic"   env:"KAFKA_TOPIC"`
		GroupId       string `yaml:"groupId"   env:"KAFKA_GROUPID"`
		NumPartitions int    `yaml:"numPartitions"   env:"KAFKA_NUMPARTITIONS"`
		Replication   int    `yaml:"replication"   env:"KAFKA_REPLICATION"`
	}

	Elasticsearch struct {
		Host    string `yaml:"host"   env:"ES_HOST"`
		PreName string `yaml:"preName"   env:"ES_HOST"`
	}

	ChunZhen struct {
		Path string `yaml:"path"   env:"CHUNZHEN_PATH"`
	}

	IP2location struct {
		Path string `yaml:"path"   env:"IP2LOCATION_PATH"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
