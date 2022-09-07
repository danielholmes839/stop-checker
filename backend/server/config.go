package server

import (
	"flag"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	OCTRANSPO_ENDPOINT  string
	OCTRANSPO_APP_ID    string
	OCTRANSPO_API_KEY   string
	GOOGLE_MAPS_API_KEY string
	SERVER_TIMEZONE     *time.Location
	SERVER_PORT         string
	SERVER_ENABLE_CORS  bool
	SERVER_DATASET      string
	DATASET_FOLDER      string
	AWS_BUCKET          string
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
	tz, _ := time.LoadLocation(viper.GetString("server.timezone"))
	return Config{
		OCTRANSPO_ENDPOINT:  viper.GetString("octranspo.endpoint"),
		OCTRANSPO_APP_ID:    viper.GetString("octranspo.app_id"),
		OCTRANSPO_API_KEY:   viper.GetString("octranspo.api_key"),
		GOOGLE_MAPS_API_KEY: viper.GetString("google_cloud.api_key"),
		SERVER_TIMEZONE:     tz,
		SERVER_PORT:         viper.GetString("server.port"),
		SERVER_ENABLE_CORS:  viper.GetBool("server.cors"),
		SERVER_DATASET:      viper.GetString("server.dataset"),
		DATASET_FOLDER:      viper.GetString("filesystem.dataset_folder"),
		AWS_BUCKET:          viper.GetString("aws.dataset_bucket"),
	}
}
