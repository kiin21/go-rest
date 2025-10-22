package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// Config Management config with Viper
type Config struct {
	ServerPort string `mapstructure:"SERVER_PORT"`
	LogLevel   string `mapstructure:"LOG_LEVEL"`

	KafkaBrokers            string `mapstructure:"KAFKA_BROKERS"`
	KafkaTopicNotifications string `mapstructure:"KAFKA_TOPIC_NOTIFICATIONS"`
	KafkaConsumerGroup      string `mapstructure:"KAFKA_CONSUMER_GROUP"`

	MongoURI        string `mapstructure:"MONGODB_URI"`
	MongoDatabase   string `mapstructure:"MONGODB_DATABASE"`
	MongoCollection string `mapstructure:"MONGODB_COLLECTION"`
}

func LoadConfig() (config Config, err error) {
	viper.SetConfigFile(".env_dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
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
