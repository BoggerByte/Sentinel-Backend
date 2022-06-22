package main

import (
	db "github.com/BoggerByte/Sentinel-backend.git/db/sqlc"
	"github.com/BoggerByte/Sentinel-backend.git/pkg"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("Failed to initialize config: %v", err.Error())
	}

	if err := initEnv(); err != nil {
		logrus.Fatalf("Failed to load env variables: %v", err.Error())
	}

	err := db.Init(db.ConnectionConfig{
		Driver:   viper.GetString("db.driver"),
		Source:   viper.GetString("db.source"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Name:     viper.GetString("db.name"),
		SSLMode:  viper.GetString("db.ssl-mode"),
	})
	if err != nil {
		logrus.Fatalf("Failed to connect to DB: %v", err.Error())
	}

	// TODo Implement DB Connection
	// TODO Implement Auth middlewares

	server := pkg.NewServer()

	if err := server.Run(viper.GetString("address")); err != nil {
		logrus.Fatalf("Error occured while running server: %v", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	return viper.ReadInConfig()
}

func initEnv() error {
	return godotenv.Load()
}
