package v1

import (
	"github.com/BoggerByte/Sentinel-backend.git/pkg/conttollers"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *conttollers.Service
}

func NewHandler(service *conttollers.Service) *Handler {
	return &Handler{}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	// controllers initialization
	router.GET("/", h.getIndex)

	api := router.Group("/api/v1")
	{
		auth := router.Group("/auth")
		{
			auth.POST("/jwt/login", h.login)
			auth.POST("/jwt/refresh", h.refresh)
		}
		servers := api.Group("/guilds/:uid")
		{
			servers.GET("", h.getGuilds)

			servers.GET("/:id", h.getGuild)

			configs := servers.Group("/:id/configs")
			{
				configs.GET("", h.getConfigs)

				configs.POST("/:v", h.createConfig)
				configs.GET("/:v", h.getConfig)
				configs.PUT("/:v", h.updateConfig)
				configs.DELETE("/:v", h.deleteConfig)
			}
		}
	}

	return router
}
