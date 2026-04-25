package config

import "github.com/spf13/viper"

type RedisConfig struct {
	Host       string
	Port       string
	Password   string
	DB         int
	MaxRetries int
	PoolSize   int
}

func loadRedisConfig() (*RedisConfig, error) {
	redisConfig := &RedisConfig{
		Host:       viper.GetString("cache.redis.host"),
		Port:       viper.GetString("cache.redis.port"),
		Password:   viper.GetString("cache.redis.password"),
		DB:         viper.GetInt("cache.redis.database"),
		MaxRetries: viper.GetInt("cache.redis.maxRetries"),
		PoolSize:   viper.GetInt("cache.redis.poolSize"),
	}

	return redisConfig, nil
}
