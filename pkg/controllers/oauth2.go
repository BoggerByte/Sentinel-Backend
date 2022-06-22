package controllers

import (
	db "github.com/BoggerByte/Sentinel-backend.git/db/sqlc"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

type Oauth2Controller struct {
	store                *db.Store
	discordOauth2Service *services.DiscordOauth2Service
}

func NewOauth2Controller(discordOauth2Service *services.DiscordOauth2Service) *Oauth2Controller {
	return &Oauth2Controller{
		store:                db.GetStore(),
		discordOauth2Service: discordOauth2Service,
	}
}

func (ctrl *Oauth2Controller) GenerateURL(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"url": ctrl.discordOauth2Service.GenerateURL("random"),
	})
}

type oauth2RedirectRequest struct {
	Code  string `form:"code" binding:"required"`
	State string `form:"state"`
}

func (ctrl *Oauth2Controller) HandleRedirect(c *gin.Context) {
	var req oauth2RedirectRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	token, err := ctrl.discordOauth2Service.Exchange(req.Code)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	dUser, err := ctrl.discordOauth2Service.GetUser(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	dGuilds, err := ctrl.discordOauth2Service.GetUserGuilds(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	arg := db.CreateAdminTxParams{
		Admin: db.CreateAdminParams{
			DiscordID:     dUser.ID,
			Username:      dUser.Username,
			Discriminator: dUser.Discriminator,
			Verified:      dUser.Verified,
			Email:         dUser.Email,
			Avatar:        dUser.Avatar,
			Flags:         dUser.Flags,
			Banner:        dUser.Banner,
			AccentColor:   dUser.AccentColor,
			PublicFlags:   dUser.PublicFlags,
		},
	}
	for _, dGuild := range dGuilds {
		if dGuild.IsOwner {
			gArg := db.CreateGuildParams{
				DiscordID: dGuild.ID,
				Name:      dGuild.Name,
				Icon:      dGuild.Icon,
				OwnerID:   dUser.ID,
			}
			arg.Guilds = append(arg.Guilds, gArg)
		}
	}

	admin, err := ctrl.store.CreateAdminTx(c, arg)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	c.JSON(http.StatusOK, admin)

	redirectURL := url.URL{Path: "/"}
	c.Redirect(http.StatusFound, redirectURL.RequestURI())
}
