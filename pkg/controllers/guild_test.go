package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/BoggerByte/Sentinel-backend.git/pkg/db/mock"
	db "github.com/BoggerByte/Sentinel-backend.git/pkg/db/sqlc"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/middlewares"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/modules/token"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/utils"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
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

func requireBodyMatchGuilds(t *testing.T, body *bytes.Buffer, guilds []db.Guild) {
	data, err := ioutil.ReadAll(body)

	require.NoError(t, err)

	var gotGuilds []db.Guild
	err = json.Unmarshal(data, &gotGuilds)

	require.NoError(t, err)
	require.Equal(t, guilds, gotGuilds)
}

func requireBodyMatchGuild(t *testing.T, body *bytes.Buffer, guild db.Guild) {
	data, err := ioutil.ReadAll(body)

	require.NoError(t, err)

	var gotGuild db.Guild
	err = json.Unmarshal(data, &gotGuild)

	require.NoError(t, err)
	require.Equal(t, guild, gotGuild)
}

func TestGuildController_Get(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	guild := generateRandomGuild()

	testCases := []struct {
		name           string
		guildDiscordID string
		buildStubs     func(store *mockdb.MockStore)
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:           "OK",
			guildDiscordID: guild.DiscordID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGuild(gomock.Any(), gomock.Eq(guild.DiscordID)).
					Times(1).
					Return(guild, nil)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				requireBodyMatchGuild(t, w.Body, guild)
			},
		},
		{
			name:           "NotFound",
			guildDiscordID: guild.DiscordID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGuild(gomock.Any(), gomock.Eq(guild.DiscordID)).
					Times(1).
					Return(db.Guild{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			name:           "InternalServerError/DBGetGuild",
			guildDiscordID: guild.DiscordID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGuild(gomock.Any(), gomock.Eq(guild.DiscordID)).
					Times(1).
					Return(db.Guild{}, sql.ErrConnDone)
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

			guildController := NewGuildController(store)
			router := gin.New()
			router.GET("/api/v1/guilds/:discord_id", guildController.Get)

			url := fmt.Sprintf("/api/v1/guilds/%s", tc.guildDiscordID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			tc.checkResponse(t, w)
		})
	}
}

func TestGuildController_GetAll(t *testing.T) {
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
				requireBodyMatchGuilds(t, w.Body, []db.Guild{guild})
			},
		},
		{
			name:             "InternalServerError/DBGetUserGuilds",
			accountDiscordID: account.DiscordID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserGuilds(gomock.Any(), gomock.Eq(account.DiscordID)).
					Times(1).
					Return([]db.GetUserGuildsRow{guildRow}, sql.ErrConnDone)
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
			router.GET("/api/v1/accounts/me/guilds", authMiddleware, guildController.GetUserAll)

			url := "/api/v1/accounts/me/guilds"
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
