package pkg

import (
	"github.com/BoggerByte/Sentinel-backend.git/pkg/controllers"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/middlewares"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/services"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	router := gin.Default()

	router.LoadHTMLGlob("./public/html/*")
	router.Static("/public", "./public")

	router.Use(middlewares.CORSMiddleware())
	router.Use(middlewares.RequestIDMiddleware())
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	api := router.Group("/api/v1")
	{
		/*** START OAUTH2 ***/
		discordOauth2Service := services.NewDiscordOauth2Service(&oauth2.Config{
			Endpoint:     discord.Endpoint,
			Scopes:       []string{discord.ScopeIdentify, discord.ScopeEmail, discord.ScopeGuilds},
			RedirectURL:  "http://localhost:8080/api/v1/oauth2/redirect",
			ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
			ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		})
		oauth := controllers.NewOauth2Controller(discordOauth2Service)

		api.GET("/oauth2/generate_url", oauth.GenerateURL)
		api.GET("/oauth2/redirect", oauth.HandleRedirect)

		/*** START AUTH ***/
		auth := new(controllers.AuthController)

		api.POST("/auth/jwt/login", auth.AuthorizeToken)
		api.POST("/auth/jwt/refresh", auth.RefreshToken)

		/*** START ADMINS ***/
		admin := controllers.NewAdminController()

		api.POST("/admins", admin.Create)
		api.GET("/admins/:id", admin.Get)

		/*** START GUILDS ***/
		guild := new(controllers.GuildController)

		api.GET("/guilds/:id", guild.Get)

		/*** START CONFIGS ***/
		config := new(controllers.ConfigController)

		api.GET("/guilds/:id/configs", config.GetAll)
		api.POST("/guilds/:id/configs/:v", config.Create)
		api.GET("/guilds/:id/configs/:v", config.Get)
		api.PUT("/guilds/:id/configs/:v", config.Update)
		api.DELETE("/guilds/:id/configs/:v", config.Delete)
	}

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Root",
		})
	})

	router.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{})
	})

	return &Server{router: router}
}

func (s *Server) Run(address string) error {
	return s.router.Run(address)
}
