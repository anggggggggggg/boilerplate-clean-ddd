package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	APP         *APPConfig
	DB          *DBConfig
	DBLog       *DBLogConfig
	DBTxRead    *DBTxReadConfig
	DBResolver  *DBResolverConfig
	DBWarehouse *DBWarehouseConfig
	Redis       *RedisConfig
	Queue       *QueueConfig
	Kafka       *KafkaConfig
	External    *ExternalConfig
}

func Load() (*Config, error) {
	// Initialize Viper
	viper.SetConfigName(".config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	appConfig, err := loadAPPConfig()
	if err != nil {
		return nil, err
	}

	dbConfig, err := loadDBConfig()
	if err != nil {
		return nil, err
	}

	dbLogConfig, err := loadDBLogConfig()
	if err != nil {
		return nil, err
	}

	dbTxReadConfig, err := loadDBTxReadConfig()
	if err != nil {
		return nil, err
	}

	dbResolverConfig, err := loadDBResolverConfig()
	if err != nil {
		return nil, err
	}

	dbWarehouseConfig, err := loadDBWarehouseConfig()
	if err != nil {
		return nil, err
	}

	redisConfig, err := loadRedisConfig()
	if err != nil {
		return nil, err
	}

	queueConfig, err := loadQueueConfig()
	if err != nil {
		return nil, err
	}

	kafkaConfig, err := loadKafkaConfig()
	if err != nil {
		return nil, err
	}

	externalConfig, err := loadExternalConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{
		APP:         appConfig,
		DB:          dbConfig,
		DBLog:       dbLogConfig,
		DBTxRead:    dbTxReadConfig,
		DBResolver:  dbResolverConfig,
		DBWarehouse: dbWarehouseConfig,
		Redis:       redisConfig,
		Queue:       queueConfig,
		Kafka:       kafkaConfig,
		External:    externalConfig,
	}

	return config, nil
}
