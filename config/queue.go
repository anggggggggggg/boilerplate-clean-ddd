package config

import "github.com/spf13/viper"

type QueueConfig struct {
	Concurrency         int
	ShutdownTimeout     int
	HealthcheckInterval int
	Redis               *QueueRedisConfig
}

type QueueRedisConfig struct {
	Host       string
	Port       string
	Password   string
	DB         int
	MaxRetries int
	PoolSize   int
}

func loadQueueConfig() (*QueueConfig, error) {
	redisConfig, err := loadQueueRedisConfig()
	if err != nil {
		return nil, err
	}

	queueConfig := &QueueConfig{
		Concurrency:         viper.GetInt("queue.option.concurrency"),
		ShutdownTimeout:     viper.GetInt("queue.option.shutdownTimeout"),
		HealthcheckInterval: viper.GetInt("queue.option.healthCheckInterval"),
		Redis:               redisConfig,
	}

	return queueConfig, nil
}

func loadQueueRedisConfig() (*QueueRedisConfig, error) {
	redisConfig := &QueueRedisConfig{
		Host:       viper.GetString("queue.redis.host"),
		Port:       viper.GetString("queue.redis.port"),
		Password:   viper.GetString("queue.redis.password"),
		DB:         viper.GetInt("queue.redis.database"),
		MaxRetries: viper.GetInt("queue.redis.maxRetries"),
		PoolSize:   viper.GetInt("queue.redis.poolSize"),
	}

	return redisConfig, nil
}
