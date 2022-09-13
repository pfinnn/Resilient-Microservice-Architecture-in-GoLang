package configuration

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	MongoConnection             string
	MongoConnectionNoColons     string
	FluentdHost                 string
	ServerAddress               string
	SmelterConnectionDefault    string
	SmelterConnection           string
	SmelterConnectionNoColons	string
	ToxiProxy_ClientConnection	string
	ToxiProxy_SmelterConnection	string
	ToxiProxy_SmelterConnectionNoColons	string
}

func ReadConfiguration(test_environment bool) Configuration {
	viper.AutomaticEnv()
	viper.SetDefault("MongoConnection", "mongodb://127.0.0.1/forge")
	viper.SetDefault("MongoConnectionNoColons", "127.0.0.1/forge")

	viper.SetDefault("FluentdHost", "localhost")

	viper.SetDefault("ServerAddress", "0.0.0.0:8080")

	viper.SetDefault("SmelterConnectionDefault", "localhost:8081")

	viper.SetDefault("SmelterConnection", "http://localhost:8081")
	viper.SetDefault("SmelterConnectionNoColons", "localhost:8081")

	viper.SetDefault("ToxiProxy_ClientConnection", "localhost:8474")

	viper.SetDefault("ToxiProxy_SmelterConnection", "http://127.0.0.1:8081")
	viper.SetDefault("ToxiProxy_SmelterConnectionNoColons", "127.0.0.1:8081")

	if test_environment {
		viper.SetDefault("SmelterConnection","http://127.0.0.1:8081")
		viper.SetDefault("SmelterConnectionNoColons", "127.0.0.1:8081")
	}


	viper.SetEnvPrefix("FORGESERVICE")

	var configuration Configuration
	configuration.MongoConnection = viper.GetString("MongoConnection")
	configuration.MongoConnectionNoColons = viper.GetString("MongoConnectionNoColons")
	configuration.FluentdHost = viper.GetString("FluentdHost")
	configuration.ServerAddress = viper.GetString("ServerAddress")
	configuration.SmelterConnectionDefault = viper.GetString("SmelterConnectionDefault")
	configuration.SmelterConnection = viper.GetString("SmelterConnection")
	configuration.SmelterConnectionNoColons = viper.GetString("SmelterConnectionNoColons")
	configuration.ToxiProxy_ClientConnection = viper.GetString("ToxiProxy_ClientConnection")
	configuration.ToxiProxy_SmelterConnection = viper.GetString("ToxiProxy_SmelterConnection")
	configuration.ToxiProxy_SmelterConnectionNoColons = viper.GetString("ToxiProxy_SmelterConnectionNoColons")

	return configuration
}
