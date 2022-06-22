package controllers

import (
	db "github.com/BoggerByte/Sentinel-backend.git/db/sqlc"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	store *db.Store
}

func (ctrl *AuthController) AuthorizeToken(c *gin.Context) {

}

func (ctrl *AuthController) RefreshToken(c *gin.Context) {

}
