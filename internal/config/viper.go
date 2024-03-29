package config

import "github.com/spf13/viper"

func NewViper(path string) *viper.Viper {
	config := viper.New()
	config.SetConfigType("json")
	config.SetConfigName("config")
	config.AddConfigPath(path)
	config.AddConfigPath("./")

	err := config.ReadInConfig()
	if err != nil {
		panic(err)
	}

	return config
}
