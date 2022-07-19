package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	memdb "github.com/BoggerByte/Sentinel-backend.git/pkg/db/memory"
	db "github.com/BoggerByte/Sentinel-backend.git/pkg/db/sqlc"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/middlewares"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/modules/token"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController struct {
	store      db.Store
	memStore   memdb.Store
	config     utils.Config
	tokenMaker token.Maker
}

func NewAuthController(
	store db.Store,
	memStore memdb.Store,
	config utils.Config,
	tokenMaker token.Maker,
) *AuthController {
	return &AuthController{
		store:      store,
		memStore:   memStore,
		config:     config,
		tokenMaker: tokenMaker,
	}
}

func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	refreshPayload := c.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

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

	newAccessToken, _, err := ctrl.tokenMaker.CreateToken(refreshPayload.UserDiscordID, ctrl.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	newRefreshToken, newRefreshPayload, err := ctrl.tokenMaker.CreateToken(refreshPayload.UserDiscordID, ctrl.config.RefreshTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	newSession, err := ctrl.store.CreateSession(c, db.CreateSessionParams{
		ID:           newRefreshPayload.ID,
		DiscordID:    newRefreshPayload.UserDiscordID,
		RefreshToken: newRefreshToken,
		UserAgent:    c.Request.UserAgent(),
		ClientIp:     c.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id":       newSession.ID,
		"access_token":     newAccessToken,
		"access_duration":  ctrl.config.AccessTokenDuration.Milliseconds(),
		"refresh_token":    newRefreshToken,
		"refresh_duration": ctrl.config.RefreshTokenDuration.Milliseconds(),
	})
}
