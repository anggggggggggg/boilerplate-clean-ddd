package config

import "github.com/spf13/viper"

type DBConfig struct {
	Host      string
	Port      string
	User      string
	Password  string
	Name      string
	Dialect   string
	SSLMode   string
	DBPooling *DBPooling
}

type DBLogConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Dialect  string
}

type DBTxReadConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Dialect  string
}

type DBResolverConfig struct {
	Sources  string
	Replicas string
}

type DBWarehouseConfig struct {
	Sources  string
	Replicas string
}

type DBPooling struct {
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
	ConnMaxIdleTime int
}

func loadDBConfig() (*DBConfig, error) {
	pooling, err := loadDBPoolingConfig()
	if err != nil {
		return nil, err
	}

	dbConfig := &DBConfig{
		Host:      viper.GetString("database.host"),
		Port:      viper.GetString("database.port"),
		User:      viper.GetString("database.user"),
		Password:  viper.GetString("database.password"),
		Name:      viper.GetString("database.name"),
		Dialect:   viper.GetString("database.dialect"),
		SSLMode:   viper.GetString("database.sslMode"),
		DBPooling: pooling,
	}

	return dbConfig, nil
}

func loadDBLogConfig() (*DBLogConfig, error) {
	dbLogConfig := &DBLogConfig{
		Host:     viper.GetString("database.log.host"),
		Port:     viper.GetString("database.log.port"),
		User:     viper.GetString("database.log.user"),
		Password: viper.GetString("database.log.password"),
		Name:     viper.GetString("database.log.name"),
		Dialect:  viper.GetString("database.log.dialect"),
	}

	return dbLogConfig, nil
}

func loadDBTxReadConfig() (*DBTxReadConfig, error) {
	dbTxReadConfig := &DBTxReadConfig{
		Host:     viper.GetString("database.resolver.sources.host"),
		Port:     viper.GetString("database.resolver.sources.port"),
		User:     viper.GetString("database.resolver.sources.user"),
		Password: viper.GetString("database.resolver.sources.password"),
		Name:     viper.GetString("database.resolver.sources.name"),
		Dialect:  viper.GetString("database.resolver.sources.dialect"),
	}

	return dbTxReadConfig, nil
}

func loadDBResolverConfig() (*DBResolverConfig, error) {
	dbResolverConfig := &DBResolverConfig{
		Sources:  viper.GetString("database.resolver.sources.name"),
		Replicas: viper.GetString("database.resolver.replicas"),
	}

	return dbResolverConfig, nil
}

func loadDBWarehouseConfig() (*DBWarehouseConfig, error) {
	dbWarehouseConfig := &DBWarehouseConfig{
		Sources:  viper.GetString("database.warehouse.host"),
		Replicas: "",
	}

	return dbWarehouseConfig, nil
}

func loadDBPoolingConfig() (*DBPooling, error) {
	pooling := &DBPooling{
		MaxIdleConns:    viper.GetInt("database.pooling.maxIdleConns"),
		MaxOpenConns:    viper.GetInt("database.pooling.maxOpenConns"),
		ConnMaxLifetime: viper.GetInt("database.pooling.connMaxLifetime"),
		ConnMaxIdleTime: viper.GetInt("database.pooling.connMaxIdleTime"),
	}

	return pooling, nil

}
