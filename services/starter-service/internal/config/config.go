package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// Config Management config with Viper
type Config struct {
	// Database
	DBURI string `mapstructure:"DB_URI"`
	// App config
	ServerPort string `mapstructure:"SERVER_PORT"`
	LogLevel   string `mapstructure:"LOG_LEVEL"`

	// Elasticsearch
	ElasticsearchAddresses string `mapstructure:"ELASTICSEARCH_ADDRESSES"`
	ElasticsearchUsername  string `mapstructure:"ELASTICSEARCH_USERNAME"`
	ElasticsearchPassword  string `mapstructure:"ELASTICSEARCH_PASSWORD"`

	// Kafka
	KafkaBrokers            string `mapstructure:"KAFKA_BROKERS"`
	KafkaTopicSyncEvents    string `mapstructure:"KAFKA_TOPIC_SYNC_EVENTS"`
	KafkaTopicNotifications string `mapstructure:"KAFKA_TOPIC_NOTIFICATIONS"`
	KafkaConsumerGroup      string `mapstructure:"KAFKA_CONSUMER_GROUP"`

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
