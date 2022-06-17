package main

import (
	"github.com/BoggerByte/Sentinel-backend.git/pkg"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/conttollers"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/repository"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/routers/v1"
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

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Failed to load env variables: %v", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("Failed to initialize db: %v", err.Error())
	}

	repositories := repository.NewRepository(db)
	services := conttollers.NewService(repositories)
	handlers := v1.NewHandler(services)
	server := new(pkg.Server)

	if err := server.Run(handlers.InitRoutes(), viper.GetString("port")); err != nil {
		logrus.Fatalf("Error occured while running http srver: %v", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")
	return viper.ReadInConfig()
}
