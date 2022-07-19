package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/BoggerByte/Sentinel-backend.git/pkg/db/mock"
	db "github.com/BoggerByte/Sentinel-backend.git/pkg/db/sqlc"
	"github.com/BoggerByte/Sentinel-backend.git/pub/objects"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGuildConfigController_Overwrite(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	guild := generateRandomGuild()

	guildConfigObj := objects.GuildConfig{
		Permissions: objects.GuildConfigPermissions{
			Edit: 40,
			Read: 0xfffffffffff,
		},
		Data: objects.GuildConfigData{
			UseConfig: false,
		},
	}
	guildConfigJSON, err := json.Marshal(guildConfigObj)
	require.NoError(t, err)
	guildConfig := db.GuildConfig{
		ID:        guild.ID,
		Json:      guildConfigJSON,
		CreatedAt: time.Time{},
	}

	testCases := []struct {
		name            string
		guildDiscordID  string
		guildConfigJSON []byte
		buildStubs      func(store *mockdb.MockStore)
		checkResponse   func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:            "OK",
			guildDiscordID:  guild.DiscordID,
			guildConfigJSON: guildConfigJSON,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateOrUpdateGuildConfig(gomock.Any(), gomock.Eq(db.CreateOrUpdateGuildConfigParams{
						DiscordID: guild.DiscordID,
						Json:      guildConfigJSON,
					})).
					Times(1).
					Return(guildConfig, nil)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:            "InternalServerError/DBCreateOrUpdateGuildConfig",
			guildDiscordID:  guild.DiscordID,
			guildConfigJSON: guildConfigJSON,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateOrUpdateGuildConfig(gomock.Any(), gomock.Eq(db.CreateOrUpdateGuildConfigParams{
						DiscordID: guild.DiscordID,
						Json:      guildConfigJSON,
					})).
					Times(1).
					Return(db.GuildConfig{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			name:            "BadRequest/JSON",
			guildDiscordID:  guild.DiscordID,
			guildConfigJSON: []byte("not_json"),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateOrUpdateGuildConfig(gomock.Any(), gomock.Any()).
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

			guildConfigController := NewGuildConfigController(store)
			router := gin.New()
			router.POST("/api/v1/guilds/:discord_id/config", guildConfigController.Overwrite)

			url := fmt.Sprintf("/api/v1/guilds/%s/config", tc.guildDiscordID)
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(tc.guildConfigJSON))
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			tc.checkResponse(t, w)
		})
	}
}

func TestGuildConfigController_Get(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	guild := generateRandomGuild()

	guildConfigObj := objects.GuildConfig{
		Permissions: objects.GuildConfigPermissions{
			Edit: 40,
			Read: 0xfffffffffff,
		},
		Data: objects.GuildConfigData{
			UseConfig: false,
		},
	}
	guildConfigJSON, err := json.Marshal(guildConfigObj)
	require.NoError(t, err)
	guildConfig := db.GuildConfig{
		ID:        guild.ID,
		Json:      guildConfigJSON,
		CreatedAt: time.Time{},
	}

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
					GetGuildConfig(gomock.Any(), gomock.Eq(guild.DiscordID)).
					Times(1).
					Return(guildConfig, nil)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:           "InternalServerError/DBGetGuildConfig",
			guildDiscordID: guild.DiscordID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGuildConfig(gomock.Any(), gomock.Eq(guild.DiscordID)).
					Times(1).
					Return(db.GuildConfig{}, sql.ErrConnDone)
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

			guildConfigController := NewGuildConfigController(store)
			router := gin.New()
			router.GET("/api/v1/guilds/:discord_id/config", guildConfigController.Get)

			url := fmt.Sprintf("/api/v1/guilds/%s/config", tc.guildDiscordID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			tc.checkResponse(t, w)
		})
	}
}
