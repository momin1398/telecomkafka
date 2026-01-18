// pkg/config/config.go
package config

import "os"

type Config struct {
	KafkaBrokers []string
	Topic        string
}

func Load() Config {
	return Config{
		KafkaBrokers: []string{os.Getenv("KAFKA_BROKER")},
		Topic:        os.Getenv("KAFKA_TOPIC"),
	}
}
