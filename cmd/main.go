package main

import (
	"github.com/BoggerByte/Sentinel-backend.git/pkg"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/controllers"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/db/sqlc"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/middlewares"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/middlewares/permissions"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/modules/token"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/services"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/util"
	_ "github.com/lib/pq"
	"github.com/ravener/discord-oauth2"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func main() {
	config, err := util.LoadConfig()
	if err != nil {
		logrus.Fatalf("Failed to initialize config: %v", err.Error())
	}

	store, err := db.Init(db.ConnectionConfig{
		Driver:   config.DBDriver,
		Protocol: config.DBProtocol,
		Username: config.DBUsername,
		Password: config.DBPassword,
		Host:     config.DBHost,
		Port:     config.DBPort,
		Name:     config.DBName,
		SSLMode:  config.DBSSLMode,
	})
	if err != nil {
		logrus.Fatalf("Failed to connect to DB: %v", err.Error())
	}

	tokenMaker, err := token.NewPasetoMaker(config.PasetoSymmetricKey)
	if err != nil {
		logrus.Fatalf("Failed to create PASeTo token maker: %v", err.Error())
	}

	discordOauth2Service := services.NewDiscordOauth2Service(&oauth2.Config{
		Endpoint:     discord.Endpoint,
		Scopes:       []string{discord.ScopeIdentify, discord.ScopeEmail, discord.ScopeGuilds},
		RedirectURL:  "http://localhost:8080/api/v1/oauth2/redirect",
		ClientID:     config.DiscordClientID,
		ClientSecret: config.DiscordClientSecret,
	})

	controllersV1 := controllers.Controllers{
		Account:     controllers.NewUserController(store),
		Auth:        controllers.NewAuthController(store, config, tokenMaker),
		Guild:       controllers.NewGuildController(store),
		GuildConfig: controllers.NewGuildConfigController(store),
		Oauth2:      controllers.NewOauth2Controller(store, config, tokenMaker, discordOauth2Service),
	}

	middlewaresV1 := middlewares.Middlewares{
		Auth: middlewares.NewAuthMiddleware(tokenMaker),
		Permissions: middlewares.Permissions{
			GuildConfig: permissions.NewGuildConfigPermissions(store),
		},
	}

	server := pkg.NewServer(controllersV1, middlewaresV1)
	if err := server.Run(config.Address); err != nil {
		logrus.Fatalf("Error occured while running server: %v", err.Error())
	}
}