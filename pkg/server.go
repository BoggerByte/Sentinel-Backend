package pkg

import (
	"github.com/BoggerByte/Sentinel-backend.git/pkg/controllers"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/middlewares"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	router *gin.Engine
}

func NewServer(controllers controllers.Controllers, middlewares middlewares.Middlewares) *Server {
	router := gin.Default()
	_ = router.SetTrustedProxies(nil)

	router.LoadHTMLGlob("./pub/html/*")
	router.Static("/pub", "./pub")

	router.Use(middlewares.CORS)
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	perms := middlewares.Permissions

	/* INITIALIZING ROUTES */

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{})
	})

	api := router.Group("/api/v1")
	{
		api.GET("/oauth2/new_url", controllers.GetNewOauth2URL)
		api.GET("/oauth2/new_invite_bot_url", controllers.GetNewInviteBotURL)
		api.GET("/oauth2/discord_callback", controllers.HandleDiscordCallback)

		api.POST("/auth/paseto/refresh", middlewares.Auth, controllers.RefreshToken)

		api.GET("/users/me", middlewares.Auth, controllers.GetUser)
		api.GET("/users/me/guilds", middlewares.Auth, controllers.GetUserGuilds)
		api.GET("/users/me/guilds/:discord_id", middlewares.Auth, controllers.GetUserGuild)

		api.GET("/guilds/configs/presets/:preset", controllers.GetGuildConfigPreset)
		api.GET("/guilds/:discord_id/config", middlewares.Auth, perms.GuildConfig.Get(), controllers.GetGuildConfig)
		api.POST("/guilds/:discord_id/config", middlewares.Auth, perms.GuildConfig.Overwrite(), controllers.OverwriteGuildConfig)
	}

	discordBotApi := router.Group("/discord_bot_api/v1").
		Use(middlewares.DiscordBotAuth)
	{
		discordBotApi.POST("/guilds", controllers.CreateOrUpdateGuilds)

		discordBotApi.GET("/guilds/configs", controllers.GetGuildsConfigs)
		discordBotApi.GET("/guilds/configs/presets/:preset", controllers.GetGuildConfigPreset)
		discordBotApi.GET("/guilds/:discord_id/config", controllers.GetGuildConfig)

		//discordBotApi.GET("/ws", controllers.Websocket)
	}

	return &Server{router: router}
}

func (s *Server) Run(address string) error {
	return s.router.Run(address)
}
