package config

import (
	"strings"

	"github.com/spf13/viper"
)

type KafkaConfig struct {
	Servers string
}

func loadKafkaConfig() (*KafkaConfig, error) {
	// Get servers as array and join them
	servers := viper.GetStringSlice("kafka.servers")
	serversStr := strings.Join(servers, ",")

	kafkaConfig := &KafkaConfig{
		Servers: serversStr,
	}

	return kafkaConfig, nil
}
