package util

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Address              string        `mapstructure:"address"`
	DBDriver             string        `mapstructure:"db-driver"`
	DBProtocol           string        `mapstructure:"db-protocol"`
	DBHost               string        `mapstructure:"db-host"`
	DBPort               string        `mapstructure:"db-port"`
	DBUsername           string        `mapstructure:"DB_USERNAME"`
	DBPassword           string        `mapstructure:"DB_PASSWORD"`
	DBName               string        `mapstructure:"db-name"`
	DBSSLMode            string        `mapstructure:"db-ssl-mode"`
	PasetoSymmetricKey   string        `mapstructure:"PASETO_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"token-access-duration"`
	RefreshTokenDuration time.Duration `mapstructure:"token-refresh-duration"`
	DiscordClientID      string        `mapstructure:"DISCORD_CLIENT_ID"`
	DiscordClientSecret  string        `mapstructure:"DISCORD_CLIENT_SECRET"`
}

func LoadConfig() (Config, error) {
	viper.AddConfigPath("./cfg/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
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
