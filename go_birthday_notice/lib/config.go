package lib

import (
	"fmt"
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigName("config")                // name of config file (without extension)
	viper.SetConfigType("toml")                  // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/config/")              // path to look for the config file in
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
