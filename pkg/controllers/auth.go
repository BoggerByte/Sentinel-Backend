package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	memdb "github.com/BoggerByte/Sentinel-backend.git/pkg/db/memory"
	db "github.com/BoggerByte/Sentinel-backend.git/pkg/db/sqlc"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/forms"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/middlewares"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/modules/token"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"net/http"
)

type AuthController struct {
	store      db.Store
	memStore   memdb.Store
	config     util.Config
	tokenMaker token.Maker
}

func NewAuthController(
	store db.Store,
	memStore memdb.Store,
	config util.Config,
	tokenMaker token.Maker,
) *AuthController {
	return &AuthController{
		store:      store,
		memStore:   memStore,
		config:     config,
		tokenMaker: tokenMaker,
	}
}

func (ctrl *AuthController) FinalizeLogin(c *gin.Context) {
	var json forms.LoginForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	oauth2Flow, err := ctrl.memStore.GetOauth2Flow(c, json.State)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err := errors.New("state not exists or expired")
			c.JSON(http.StatusMethodNotAllowed, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if !oauth2Flow.Completed {
		err := errors.New("oauth2 flow is not completed")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, _, err := ctrl.tokenMaker.CreateToken(oauth2Flow.UserDiscordID, ctrl.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	refreshToken, refreshPayload, err := ctrl.tokenMaker.CreateToken(oauth2Flow.UserDiscordID, ctrl.config.RefreshTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := ctrl.store.CreateSession(c, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		DiscordID:    refreshPayload.UserDiscordID,
		RefreshToken: refreshToken,
		UserAgent:    c.Request.UserAgent(),
		ClientIp:     c.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = ctrl.memStore.DeleteOauth2Flow(c, json.State)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id":    session.ID,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
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

	accessToken, _, err := ctrl.tokenMaker.CreateToken(refreshPayload.UserDiscordID, ctrl.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}
