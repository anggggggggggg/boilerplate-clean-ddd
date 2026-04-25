package config

import "github.com/spf13/viper"

type ExternalConfig struct {
	Host          string
	Token         string
	SeamlessHost  string
	SeamlessToken string
}

func loadExternalConfig() (*ExternalConfig, error) {
	externalConfig := &ExternalConfig{
		Host:          viper.GetString("external.host.reportto"),
		Token:         viper.GetString("external.token.reportto"),
		SeamlessHost:  viper.GetString("external.seamlesshost.wallet"),
		SeamlessToken: viper.GetString("external.seamlesstoken.wallet"),
	}

	return externalConfig, nil
}
