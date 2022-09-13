package configuration

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	MySQLConnectionWithDatabase string
	MySQLConnection             string
	MySQLDatabase               string
	FluentdHost    				string
	ServerAddress               string
}

func ReadConfiguration() Configuration {
	viper.AutomaticEnv()

	viper.SetDefault("MySQLConnection", "root:root@/")
	viper.SetDefault("MySQLDatabase", "smelter")
	viper.SetDefault("FluentdHost", "localhost")
	viper.SetDefault("ServerAddress", "0.0.0.0:8081")

	viper.SetEnvPrefix("SMELTERSERVICE")

	var configuration Configuration
	configuration.MySQLConnection = viper.GetString("MySQLConnection")
	configuration.MySQLDatabase = viper.GetString("MySQLDatabase")
	configuration.FluentdHost = viper.GetString("FluentdHost")
	configuration.ServerAddress = viper.GetString("ServerAddress")

	configuration.MySQLConnectionWithDatabase = configuration.MySQLConnection + configuration.MySQLDatabase

	return configuration
}
