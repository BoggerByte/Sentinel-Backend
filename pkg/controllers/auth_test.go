package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	memdb "github.com/BoggerByte/Sentinel-backend.git/pkg/db/memory"
	mockmemdb "github.com/BoggerByte/Sentinel-backend.git/pkg/db/memory_mock"
	mockdb "github.com/BoggerByte/Sentinel-backend.git/pkg/db/mock"
	db "github.com/BoggerByte/Sentinel-backend.git/pkg/db/sqlc"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/middlewares"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/modules/token"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/util"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func generateRandomSession(payload *token.Payload) db.Session {
	return db.Session{
		ID:           payload.ID,
		DiscordID:    payload.UserDiscordID,
		RefreshToken: "",
		UserAgent:    gofakeit.UserAgent(),
		ClientIp:     gofakeit.IPv4Address(),
		IsBlocked:    false,
		ExpiresAt:    payload.ExpiredAt,
		CreatedAt:    gofakeit.Date(),
	}
}

func TestAuthController_FinalizeLogin(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	config := util.Config{
		AccessTokenDuration:     15 * time.Minute,
		RefreshTokenDuration:    2 * time.Hour,
		Oauth2FlowStateDuration: time.Hour,
	}

	state := util.RandomString(32)
	stateJSON, err := json.Marshal(struct {
		State string `json:"state"`
	}{
		State: state,
	})

	tokenMaker, _ := token.NewPasetoMaker(util.RandomString(32))
	userDiscordID := util.RandomSnowflakeID().Int64()
	_, payload, err := tokenMaker.CreateToken(userDiscordID, 5*time.Minute)
	require.NoError(t, err)
	session := generateRandomSession(payload)

	testCases := []struct {
		name          string
		stateJSON     []byte
		buildStubs    func(store *mockdb.MockStore, memStore *mockmemdb.MockStore)
		checkResponse func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			stateJSON: stateJSON,
			buildStubs: func(store *mockdb.MockStore, memStore *mockmemdb.MockStore) {
				memStore.EXPECT().
					GetOauth2Flow(gomock.Any(), gomock.Eq(state)).
					Times(1).
					Return(memdb.Oauth2Flow{
						Completed:     true,
						UserDiscordID: userDiscordID,
					}, nil)
				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(session, nil)
				memStore.EXPECT().
					DeleteOauth2Flow(gomock.Any(), gomock.Eq(state)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:      "StateNotExistsOrExpired",
			stateJSON: stateJSON,
			buildStubs: func(store *mockdb.MockStore, memStore *mockmemdb.MockStore) {
				memStore.EXPECT().
					GetOauth2Flow(gomock.Any(), gomock.Eq(state)).
					Times(1).
					Return(memdb.Oauth2Flow{}, redis.Nil)
				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(0)
				memStore.EXPECT().
					DeleteOauth2Flow(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusMethodNotAllowed, w.Code)
			},
		},
		{
			name:      "Oauth2FlowIsNotCompleted",
			stateJSON: stateJSON,
			buildStubs: func(store *mockdb.MockStore, memStore *mockmemdb.MockStore) {
				memStore.EXPECT().
					GetOauth2Flow(gomock.Any(), gomock.Eq(state)).
					Times(1).
					Return(memdb.Oauth2Flow{
						Completed: false,
					}, nil)
				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(0)
				memStore.EXPECT().
					DeleteOauth2Flow(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			name:      "InternalServerError/GetOauth2Flow",
			stateJSON: stateJSON,
			buildStubs: func(store *mockdb.MockStore, memStore *mockmemdb.MockStore) {
				memStore.EXPECT().
					GetOauth2Flow(gomock.Any(), gomock.Eq(state)).
					Times(1).
					Return(memdb.Oauth2Flow{}, redis.ErrClosed)
				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(0)
				memStore.EXPECT().
					DeleteOauth2Flow(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			name:      "InternalServerError/CreateSession",
			stateJSON: stateJSON,
			buildStubs: func(store *mockdb.MockStore, memStore *mockmemdb.MockStore) {
				memStore.EXPECT().
					GetOauth2Flow(gomock.Any(), gomock.Eq(state)).
					Times(1).
					Return(memdb.Oauth2Flow{
						Completed:     true,
						UserDiscordID: userDiscordID,
					}, nil)
				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Session{}, sql.ErrConnDone)
				memStore.EXPECT().
					DeleteOauth2Flow(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			name:      "InternalServerError/DeleteOauth2Flow",
			stateJSON: stateJSON,
			buildStubs: func(store *mockdb.MockStore, memStore *mockmemdb.MockStore) {
				memStore.EXPECT().
					GetOauth2Flow(gomock.Any(), gomock.Eq(state)).
					Times(1).
					Return(memdb.Oauth2Flow{
						Completed:     true,
						UserDiscordID: userDiscordID,
					}, nil)
				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(session, nil)
				memStore.EXPECT().
					DeleteOauth2Flow(gomock.Any(), gomock.Eq(state)).
					Times(1).
					Return(redis.ErrClosed)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			name:      "BadRequest",
			stateJSON: []byte("not_json"),
			buildStubs: func(store *mockdb.MockStore, memStore *mockmemdb.MockStore) {
				memStore.EXPECT().
					GetOauth2Flow(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(0)
				memStore.EXPECT().
					DeleteOauth2Flow(gomock.Any(), gomock.Any()).
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
			memStore := mockmemdb.NewMockStore(ctrl)
			tc.buildStubs(store, memStore)

			router := gin.New()
			authController := NewAuthController(store, memStore, config, tokenMaker)
			router.POST("/finalize_login", authController.FinalizeLogin)

			req, err := http.NewRequest(http.MethodPost, "/finalize_login", bytes.NewBuffer(tc.stateJSON))
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			tc.checkResponse(t, w)
		})
	}
}

func TestAuthController_RefreshToken(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	config := util.Config{
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 2 * time.Hour,
	}

	tokenMaker, _ := token.NewPasetoMaker(util.RandomString(32))
	refreshToken, refreshPayload, err := tokenMaker.CreateToken(1234, time.Minute)
	require.NoError(t, err)
	session := generateRandomSession(refreshPayload)

	testCases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(session, nil)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name: "SessionNotFound",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(db.Session{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			name: "Unauthorized/SessionBlocked",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(db.Session{
						IsBlocked: true,
					}, nil)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			name: "Unauthorized/DiscordIDMismatch",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(db.Session{
						IsBlocked: false,
						DiscordID: -1,
					}, nil)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			name: "InternalServerError/DBGetSession",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(db.Session{}, sql.ErrConnDone)
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

			router := gin.New()
			authMiddleware := middlewares.NewAuthMiddleware(tokenMaker)
			authController := NewAuthController(store, nil, config, tokenMaker)
			router.GET("/refresh", authMiddleware, authController.RefreshToken)

			req, err := http.NewRequest(http.MethodGet, "/refresh", nil)
			require.NoError(t, err)
			authHeader := fmt.Sprintf("%s %s", middlewares.AuthorizationTypeBearer, refreshToken)
			req.Header.Set(middlewares.AuthorizationHeaderKey, authHeader)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			tc.checkResponse(t, w)
		})
	}
}
