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
		api.GET("/oauth2/new_url", controllers.Oauth2.NewURL)
		api.GET("/oauth2/discord_callback", controllers.Oauth2.DiscordCallback)

		api.POST("/auth/paseto/finalize_login", controllers.Auth.FinalizeLogin)
		api.POST("/auth/paseto/refresh", middlewares.Auth, controllers.Auth.RefreshToken)

		api.GET("/users/me", middlewares.Auth, controllers.Account.Get)
		api.GET("/users/me/guilds", middlewares.Auth, controllers.Guild.GetAll)

		api.GET("/guilds/:discord_id", controllers.Guild.Get)

		api.POST("/guilds/:discord_id/config", middlewares.Auth, perms.GuildConfig.Overwrite(), controllers.GuildConfig.Overwrite)
		api.GET("/guilds/:discord_id/config", middlewares.Auth, perms.GuildConfig.Get(), controllers.GuildConfig.Get)
	}

	return &Server{router: router}
}

func (s *Server) Run(address string) error {
	return s.router.Run(address)
}
