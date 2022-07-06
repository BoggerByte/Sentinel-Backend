package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	db "github.com/BoggerByte/Sentinel-backend.git/pkg/db/sqlc"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/middlewares"
	token2 "github.com/BoggerByte/Sentinel-backend.git/pkg/modules/token"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController struct {
	store      db.Store
	config     util.Config
	tokenMaker token2.Maker
}

func NewAuthController(store db.Store, config util.Config, tokenMaker token2.Maker) *AuthController {
	return &AuthController{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
	}
}

func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	refreshPayload := c.MustGet(middlewares.AuthorizationPayloadKey).(*token2.Payload)

	session, err := ctrl.store.GetSession(c, refreshPayload.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.DiscordID != refreshPayload.UserDiscordID {
		err := fmt.Errorf("incorrect session user")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, _, err := ctrl.tokenMaker.CreateToken(refreshPayload.UserDiscordID, ctrl.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}
