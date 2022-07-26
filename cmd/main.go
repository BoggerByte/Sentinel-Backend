package main

import (
	"github.com/BoggerByte/Sentinel-backend.git/pkg"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/controllers"
	memdb "github.com/BoggerByte/Sentinel-backend.git/pkg/db/memory"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/db/sqlc"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/middlewares"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/middlewares/permissions"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/modules/token"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/services"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/utils"
	"github.com/gin-contrib/cors"
	_ "github.com/lib/pq"
	"github.com/ravener/discord-oauth2"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"time"
)

func main() {
	config, err := utils.LoadConfig()
	if err != nil {
		logrus.Fatalf("Failed to initialize config: %v", err.Error())
	}

	store, err := db.NewSQLStore(db.ConnectionConfig{
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

	memStore, err := memdb.NewRedisStore(memdb.ConnectionConfig{
		Host:     config.RedisHost,
		Port:     config.RedisPort,
		Password: config.RedisPassword,
		DB:       0, // default
	})
	if err != nil {
		logrus.Fatalf("Failed to connect to Memomry DB: %s", err.Error())
	}

	tokenMaker, err := token.NewPasetoMaker(config.PasetoSymmetricKey)
	if err != nil {
		logrus.Fatalf("Failed to create PASeTo token maker: %v", err.Error())
	}

	discordOauth2Service := services.NewDiscordOauth2Service(&oauth2.Config{
		Endpoint:     discord.Endpoint,
		Scopes:       []string{discord.ScopeIdentify, discord.ScopeEmail, discord.ScopeGuilds},
		RedirectURL:  "http://localhost:5173/oauth2/discord_callback",
		ClientID:     config.DiscordClientID,
		ClientSecret: config.DiscordClientSecret,
	})

	controllersV1 := controllers.Controllers{
		User:        controllers.NewUserController(store),
		Auth:        controllers.NewAuthController(store, memStore, config, tokenMaker),
		Guild:       controllers.NewGuildController(store),
		GuildConfig: controllers.NewGuildConfigController(store),
		Oauth2:      controllers.NewOauth2Controller(store, memStore, config, tokenMaker, discordOauth2Service),
	}
	middlewaresV1 := middlewares.Middlewares{
		CORS: cors.New(cors.Config{
			AllowAllOrigins:        true,
			AllowMethods:           []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders:           []string{"Content-Type", "Origin", "Access-Control-Allow-Origin", "Authorization", "Accept", "Accept-Encoding"},
			AllowCredentials:       true,
			ExposeHeaders:          []string{"Content-Length"},
			MaxAge:                 12 * time.Hour,
			AllowBrowserExtensions: true,
			AllowWebSockets:        true,
			AllowFiles:             true,
		}),
		Auth:           middlewares.NewAuthMiddleware(tokenMaker),
		DiscordBotAuth: middlewares.NewDiscordBotAuthMiddleware(config),
		Permissions: middlewares.Permissions{
			GuildConfig: permissions.NewGuildConfigPermissions(store),
		},
	}

	server := pkg.NewServer(controllersV1, middlewaresV1)
	if err := server.Run(config.ServerHTTPAddress); err != nil {
		logrus.Fatalf("Error occured while running server: %v", err.Error())
	}
}
