package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// Config Management config with Viper
type Config struct {
	// App config
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPass     string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	ServerPort string `mapstructure:"SERVER_PORT"`
	LogLevel   string `mapstructure:"LOG_LEVEL"`

	// Elasticsearch
	ElasticsearchAddresses string `mapstructure:"ELASTICSEARCH_ADDRESSES"`
	ElasticsearchUsername  string `mapstructure:"ELASTICSEARCH_USERNAME"`
	ElasticsearchPassword  string `mapstructure:"ELASTICSEARCH_PASSWORD"`

	// Kafka
	KafkaBrokers         string `mapstructure:"KAFKA_BROKERS"`
	KafkaTopicSyncEvents string `mapstructure:"KAFKA_TOPIC_SYNC_EVENTS"`
	KafkaConsumerGroup   string `mapstructure:"KAFKA_CONSUMER_GROUP"`

	// HTTP metadata
	PublicBaseURL string `mapstructure:"PUBLIC_BASE_URL"`

	// DBDriver      string `mapstructure:"DB_DRIVER"`
	// AppVersion    string `mapstructure:"APP_VERSION"`
	// ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig() (config Config, err error) {
	viper.SetConfigFile(".env_dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		// If the config file is not found, return a specific error
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return config, fmt.Errorf("config file not found: %w", err)
		}
		return config, fmt.Errorf("error reading config file: %w", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("unable to decode config into struct: %w", err)
	}
	return
}
