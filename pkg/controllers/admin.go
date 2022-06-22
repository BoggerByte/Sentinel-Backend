package controllers

import (
	"database/sql"
	db "github.com/BoggerByte/Sentinel-backend.git/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AdminController struct {
	store *db.Store
}

func NewAdminController() *AdminController {
	return &AdminController{store: db.GetStore()}
}

type createAdminRequest struct {
	DiscordID     string `json:"discord_id" binding:"required"`
	Username      string `json:"username" binding:"required"`
	Discriminator string `json:"discriminator" binding:"required"`
	Verified      bool   `json:"verified" binding:"required"`
	Email         string `json:"email" binding:"required"`
	Avatar        string `json:"avatar"`
	Flags         int64  `json:"flags"`
	Banner        string `json:"banner"`
	AccentColor   int64  `json:"accent_color"`
	PublicFlags   int64  `json:"public_flags"`
}

func (ctrl *AdminController) Create(c *gin.Context) {
	var req createAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	arg := db.CreateAdminParams{
		DiscordID:     req.DiscordID,
		Username:      req.Username,
		Discriminator: req.Discriminator,
		Verified:      req.Verified,
		Email:         req.Email,
		Avatar:        req.Avatar,
		Flags:         req.Flags,
		Banner:        req.Banner,
		AccentColor:   req.AccentColor,
		PublicFlags:   req.PublicFlags,
	}
	admin, err := ctrl.store.CreateAdmin(c, arg)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, admin)
}

type getAdminRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (ctrl *AdminController) Get(c *gin.Context) {
	var req getAdminRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	admin, err := ctrl.store.GetAdmin(c, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, admin)
}
