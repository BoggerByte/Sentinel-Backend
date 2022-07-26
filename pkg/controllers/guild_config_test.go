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

func TestGuildConfigController_GetGuildConfigPreset(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	testCases := []struct {
		name          string
		presetName    string
		checkResponse func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:       "OK/default",
			presetName: "default",
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:       "BadRequest",
			presetName: "invalid",
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			guildConfigController := NewGuildConfigController(nil)
			router := gin.New()
			router.GET("/api/v1/guilds/configs/presets/:preset", guildConfigController.GetGuildConfigPreset)

			url := fmt.Sprintf("/api/v1/guilds/configs/presets/%s", tc.presetName)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			tc.checkResponse(t, w)
		})
	}
}

func TestGuildConfigController_OverwriteGuildConfig(t *testing.T) {
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
		Preset: "default",
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
			router.POST("/api/v1/guilds/:discord_id/config", guildConfigController.OverwriteGuildConfig)

			url := fmt.Sprintf("/api/v1/guilds/%s/config", tc.guildDiscordID)
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(tc.guildConfigJSON))
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			tc.checkResponse(t, w)
		})
	}
}

func TestGuildConfigController_GetGuildConfig(t *testing.T) {
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
		Preset: "default",
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
			router.GET("/api/v1/guilds/:discord_id/config", guildConfigController.GetGuildConfig)

			url := fmt.Sprintf("/api/v1/guilds/%s/config", tc.guildDiscordID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			tc.checkResponse(t, w)
		})
	}
}

func TestGuildConfigController_GetGuildsConfigs(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	guild := generateRandomGuild()
	guildConfigObj := objects.DefaultGuildConfig
	guildConfigJSON, err := json.Marshal(guildConfigObj)
	require.NoError(t, err)
	guildConfig := db.GuildConfig{
		ID:        guild.ID,
		Json:      guildConfigJSON,
		CreatedAt: time.Time{},
	}

	testCases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGuildsConfigs(gomock.Any()).
					Times(1).
					Return([]db.GuildConfig{guildConfig}, nil)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name: "InternalServerError/GetGuildsConfigs",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGuildsConfigs(gomock.Any()).
					Times(1).
					Return([]db.GuildConfig{}, sql.ErrNoRows)
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
			router.GET("/api/v1/guilds/configs", guildConfigController.GetGuildsConfigs)

			url := fmt.Sprintf("/api/v1/guilds/configs")
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			tc.checkResponse(t, w)
		})
	}
}
