package application

import (
	"flag"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	OCTRANSPO_ENDPOINT       string
	OCTRANSPO_APP_ID         string
	OCTRANSPO_API_KEY        string
	GOOGLE_MAPS_API_KEY      string
	SERVER_PORT              string
	SERVER_ENABLE_CORS       bool
	SERVER_ENABLE_PLAYGROUND bool
	DATA_GTFS                string
	DATA_DIRECTIONS          string
	OSRM_ENDPOINT            string
}

func GetConfig() Config {
	return Config{
		OCTRANSPO_ENDPOINT:       viper.GetString("octranspo.endpoint"),
		OCTRANSPO_APP_ID:         viper.GetString("octranspo.app_id"),
		OCTRANSPO_API_KEY:        viper.GetString("octranspo.api_key"),
		GOOGLE_MAPS_API_KEY:      viper.GetString("google_cloud.api_key"),
		SERVER_PORT:              viper.GetString("server.port"),
		SERVER_ENABLE_CORS:       viper.GetBool("server.cors"),
		SERVER_ENABLE_PLAYGROUND: viper.GetBool("server.playground"),
		DATA_GTFS:                viper.GetString("data.gtfs"),
		DATA_DIRECTIONS:          viper.GetString("data.directions"),
		OSRM_ENDPOINT:            viper.GetString("osrm.endpoint"),
	}
}

func ReadConfig() {
	// command line setup
	flag.String("config", "dev", "run with '--config=dev' to use the dev.toml file")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// config setup
	config := viper.GetString("config")
	viper.SetConfigName(config)
	viper.SetConfigType("toml")
	viper.AddConfigPath("./")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
