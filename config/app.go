package config

import "github.com/spf13/viper"

type APPConfig struct {
	Name        string
	Desc        string
	Mode        string
	Port        string
	Env         string
	Prefork     string
	AUWI        string
	ProgramFile string
	WorkingDir  string
}

func loadAPPConfig() (*APPConfig, error) {
	appConfig := &APPConfig{
		Name:        viper.GetString("app.name"),
		Desc:        viper.GetString("app.desc"),
		Mode:        viper.GetString("app.mode"),
		Port:        viper.GetString("app.port"),
		Env:         viper.GetString("app.env"),
		Prefork:     viper.GetString("app.prefork"),
		AUWI:        viper.GetString("app.auwi"),
		ProgramFile: viper.GetString("app.program_file"),
		WorkingDir:  viper.GetString("app.working_dir"),
	}

	return appConfig, nil
}
