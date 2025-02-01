package config

import (
	"github.com/prajwalpamin/banking-ledger/pkg/database"
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort     int    `mapstructure:"SERVER_PORT"`
	Environment    string `mapstructure:"ENVIRONMENT"`
	PostgresConfig database.PostgresConfig
	MongoConfig    database.MongoConfig
	KafkaConfig    KafkaConfig
}

type KafkaConfig struct {
	Brokers []string `mapstructure:"KAFKA_BROKERS"`
}

func Load() (*Config, error) {
	viper.SetDefault("SERVER_PORT", 8080)
	viper.SetDefault("ENVIRONMENT", "development")

	viper.AutomaticEnv()

	cfg := &Config{
		ServerPort:  viper.GetInt("SERVER_PORT"),
		Environment: viper.GetString("ENVIRONMENT"),
		PostgresConfig: database.PostgresConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetInt("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		MongoConfig: database.MongoConfig{
			URI:      viper.GetString("MONGODB_URI"),
			Database: viper.GetString("MONGODB_DATABASE"),
		},
		KafkaConfig: KafkaConfig{
			Brokers: viper.GetStringSlice("KAFKA_BROKERS"),
		},
	}

	return cfg, nil
}
