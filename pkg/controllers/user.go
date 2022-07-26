package controllers

import (
	"database/sql"
	db "github.com/BoggerByte/Sentinel-backend.git/pkg/db/sqlc"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/middlewares"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/modules/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	store db.Store
}

func NewUserController(store db.Store) *UserController {
	return &UserController{store: store}
}

func (ctrl *UserController) GetUser(c *gin.Context) {
	payload := c.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

	account, err := ctrl.store.GetUser(c, payload.UserDiscordID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, account)
}
