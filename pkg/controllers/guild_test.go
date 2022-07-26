package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/BoggerByte/Sentinel-backend.git/pkg/db/mock"
	db "github.com/BoggerByte/Sentinel-backend.git/pkg/db/sqlc"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/forms"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/middlewares"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/modules/token"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/utils"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func generateRandomGuild() db.Guild {
	return db.Guild{
		ID:             int64(utils.RandomInt(1, 1000)),
		DiscordID:      utils.RandomSnowflakeID().String(),
		Name:           gofakeit.AppName(),
		Icon:           gofakeit.ImageURL(400, 400),
		OwnerDiscordID: utils.RandomSnowflakeID().String(),
	}
}

func TestGuildController_GetUserGuild(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	guild := generateRandomGuild()
	account := generateRandomUser()
	accountGuildRel := db.UserGuild{
		GuildDiscordID:   guild.DiscordID,
		AccountDiscordID: account.DiscordID,
		Permissions:      1 << utils.RandomInt(0, 40),
	}
	guildRow := db.GetUserGuildRow{
		ID:             guild.ID,
		DiscordID:      guild.DiscordID,
		OwnerDiscordID: guild.OwnerDiscordID,
		Icon:           guild.Icon,
		Name:           guild.Name,
		Permissions:    accountGuildRel.Permissions,
		ConfigRead:     0xfffffffff,
		ConfigEdit:     40,
	}

	testCases := []struct {
		name             string
		accountDiscordID string
		buildStubs       func(store *mockdb.MockStore)
		checkResponse    func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:             "OK",
			accountDiscordID: account.DiscordID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserGuild(gomock.Any(), gomock.Eq(db.GetUserGuildParams{
						AccountDiscordID: account.DiscordID,
						GuildDiscordID:   guild.DiscordID,
					})).
					Times(1).
					Return(guildRow, nil)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:             "NotFound",
			accountDiscordID: account.DiscordID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserGuild(gomock.Any(), gomock.Eq(db.GetUserGuildParams{
						AccountDiscordID: account.DiscordID,
						GuildDiscordID:   guild.DiscordID,
					})).
					Times(1).
					Return(db.GetUserGuildRow{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			name:             "InternalServerError/GetUserGuild",
			accountDiscordID: account.DiscordID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserGuild(gomock.Any(), gomock.Eq(db.GetUserGuildParams{
						AccountDiscordID: account.DiscordID,
						GuildDiscordID:   guild.DiscordID,
					})).
					Times(1).
					Return(db.GetUserGuildRow{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			tokenMaker, _ := token.NewPasetoMaker(utils.RandomString(32))
			authMiddleware := middlewares.NewAuthMiddleware(tokenMaker)
			guildController := NewGuildController(store)
			router := gin.New()
			router.GET("/api/v1/users/me/guilds/:discord_id", authMiddleware, guildController.GetUserGuild)

			url := fmt.Sprintf("/api/v1/users/me/guilds/%s", guild.DiscordID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			accessToken, _, err := tokenMaker.CreateToken(tc.accountDiscordID, time.Minute)
			require.NoError(t, err)

			authHeader := fmt.Sprintf("%s %s", middlewares.AuthorizationTypeBearer, accessToken)
			req.Header.Set(middlewares.AuthorizationHeaderKey, authHeader)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			tc.checkResponse(t, w)
		})
	}
}

func TestGuildController_GetUserGuilds(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	guild := generateRandomGuild()
	account := generateRandomUser()
	accountGuildRel := db.UserGuild{
		GuildDiscordID:   guild.DiscordID,
		AccountDiscordID: account.DiscordID,
		Permissions:      1 << utils.RandomInt(0, 40),
	}
	guildRow := db.GetUserGuildsRow{
		ID:             guild.ID,
		DiscordID:      guild.DiscordID,
		OwnerDiscordID: guild.OwnerDiscordID,
		Icon:           guild.Icon,
		Name:           guild.Name,
		Permissions:    accountGuildRel.Permissions,
		ConfigRead:     0xfffffffff,
		ConfigEdit:     40,
	}

	testCases := []struct {
		name             string
		accountDiscordID string
		buildStubs       func(store *mockdb.MockStore)
		checkResponse    func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:             "OK",
			accountDiscordID: account.DiscordID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserGuilds(gomock.Any(), gomock.Eq(account.DiscordID)).
					Times(1).
					Return([]db.GetUserGuildsRow{guildRow}, nil)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:             "InternalServerError/DBGetUserGuilds",
			accountDiscordID: account.DiscordID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserGuilds(gomock.Any(), gomock.Eq(account.DiscordID)).
					Times(1).
					Return([]db.GetUserGuildsRow{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			tokenMaker, _ := token.NewPasetoMaker(utils.RandomString(32))
			authMiddleware := middlewares.NewAuthMiddleware(tokenMaker)
			guildController := NewGuildController(store)
			router := gin.New()
			router.GET("/api/v1/users/me/guilds", authMiddleware, guildController.GetUserGuilds)

			url := "/api/v1/users/me/guilds"
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			accessToken, _, err := tokenMaker.CreateToken(tc.accountDiscordID, time.Minute)
			require.NoError(t, err)

			authHeader := fmt.Sprintf("%s %s", middlewares.AuthorizationTypeBearer, accessToken)
			req.Header.Set(middlewares.AuthorizationHeaderKey, authHeader)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			tc.checkResponse(t, w)
		})
	}
}

func TestGuildController_CreateOrUpdateGuild(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	guild := generateRandomGuild()
	guildJSON, err := json.Marshal(forms.CreateOrUpdateGuildsJSON{
		Guilds: []forms.CreateOrUpdateGuildJSON{
			{
				DiscordID:      guild.DiscordID,
				Name:           guild.Name,
				Icon:           guild.Icon,
				OwnerDiscordID: guild.OwnerDiscordID,
			},
		}})
	require.NoError(t, err)

	testCases := []struct {
		name          string
		guildJSON     []byte
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			guildJSON: guildJSON,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ExecTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:      "InternalServerError/ExecTx",
			guildJSON: guildJSON,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ExecTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			name:      "BadRequest/JSON",
			guildJSON: []byte("invalid JSON"),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ExecTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			guildController := NewGuildController(store)
			router := gin.New()
			router.POST("/api/v1/guilds", guildController.CreateOrUpdateGuilds)

			url := "/api/v1/guilds"
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(tc.guildJSON))
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			tc.checkResponse(t, w)
		})
	}
}

func TestGuildController_CreateOrUpdateGuilds(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
}
