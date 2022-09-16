package server

import (
	"flag"

	"github.com/rs/zerolog/log"
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
	SERVER_DATASET           string
}

func ReadConfig() error {
	viper.SetDefault("config", "dev")

	// read the command line
	flag.String("config", "dev", "set the config name default is 'dev'")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	config := viper.GetString("config")
	log.Info().Str("config", config).Msg("read config")

	viper.SetConfigName(config)
	viper.SetConfigType("toml")
	viper.AddConfigPath("./cmd/server")

	return viper.ReadInConfig()
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
		SERVER_DATASET:           viper.GetString("server.dataset"),
	}
}
