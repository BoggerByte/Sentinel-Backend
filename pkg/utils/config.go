package utils

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	ServerHTTPAddress           string        `mapstructure:"SERVER_HTTP_ADDRESS"`
	DBDriver                    string        `mapstructure:"DB_DRIVER"`
	DBProtocol                  string        `mapstructure:"DB_PROTOCOL"`
	DBHost                      string        `mapstructure:"DB_HOST"`
	DBPort                      string        `mapstructure:"DB_PORT"`
	DBUsername                  string        `mapstructure:"DB_USERNAME"`
	DBPassword                  string        `mapstructure:"DB_PASSWORD"`
	DBName                      string        `mapstructure:"DB_NAME"`
	DBSSLMode                   string        `mapstructure:"DB_SSL_MODE"`
	RedisHost                   string        `mapstructure:"REDIS_HOST"`
	RedisPort                   string        `mapstructure:"REDIS_PORT"`
	RedisPassword               string        `mapstructure:"REDIS_PASSWORD"`
	PasetoSymmetricKey          string        `mapstructure:"PASETO_SYMMETRIC_KEY"`
	AccessTokenDuration         time.Duration `mapstructure:"TOKEN_ACCESS_DURATION"`
	RefreshTokenDuration        time.Duration `mapstructure:"TOKEN_REFRESH_DURATION"`
	Oauth2FlowStateDuration     time.Duration `mapstructure:"OAUTH2_FLOW_STATE_DURATION"`
	DiscordClientID             string        `mapstructure:"DISCORD_CLIENT_ID"`
	DiscordClientSecret         string        `mapstructure:"DISCORD_CLIENT_SECRET"`
	DiscordBotSentinelAPISecret string        `mapstructure:"DISCORD_BOT_SENTINEL_API_SECRET"`
}

func LoadConfig() (Config, error) {
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	viper.AddConfigPath("./")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	if err := viper.MergeInConfig(); err != nil {
		return Config{}, err
	}

	viper.AutomaticEnv()

	var config Config
	err := viper.Unmarshal(&config)
	return config, err
}
