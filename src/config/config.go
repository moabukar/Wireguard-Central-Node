package config

import (
	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.SetConfigName("config")              // name of config file (without extension)
	viper.SetConfigType("json")                // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/opt/wgManagerAuth/") // path to look for the config file in

	//DEFAULTS
	viper.SetDefault("SERVER.PORT", "8443")   //PORT
	viper.SetDefault("SERVER.SECURITY", true) //SECURITY

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic("Cannot read env file!")
	}
}
