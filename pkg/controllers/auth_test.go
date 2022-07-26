package controllers

import (
	"database/sql"
	"fmt"
	memdb "github.com/BoggerByte/Sentinel-backend.git/pkg/db/memory"
	mockmemdb "github.com/BoggerByte/Sentinel-backend.git/pkg/db/memory_mock"
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

func generateRandomSession(payload *token.Payload) memdb.Session {
	return memdb.Session{
		ID:           payload.ID,
		DiscordID:    payload.UserDiscordID,
		RefreshToken: "",
		UserAgent:    gofakeit.UserAgent(),
		ClientIp:     gofakeit.IPv4Address(),
		IsBlocked:    false,
	}
}

func TestAuthController_RefreshToken(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	config := utils.Config{
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 2 * time.Hour,
	}

	tokenMaker, _ := token.NewPasetoMaker(utils.RandomString(32))
	refreshToken, refreshPayload, err := tokenMaker.CreateToken("1234", time.Minute)
	require.NoError(t, err)
	session := generateRandomSession(refreshPayload)

	testCases := []struct {
		name          string
		buildStubs    func(store *mockmemdb.MockStore)
		checkResponse func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(store *mockmemdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(session, nil)
				store.EXPECT().
					SetSession(gomock.Any(), gomock.Any(), gomock.Eq(config.RefreshTokenDuration)).
					Times(1).
					Return(session, nil)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name: "SessionNotFound",
			buildStubs: func(store *mockmemdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(memdb.Session{}, sql.ErrNoRows)
				store.EXPECT().
					SetSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			name: "Unauthorized/SessionBlocked",
			buildStubs: func(store *mockmemdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(memdb.Session{
						IsBlocked: true,
					}, nil)
				store.EXPECT().
					SetSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			name: "Unauthorized/DiscordIDMismatch",
			buildStubs: func(store *mockmemdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(memdb.Session{
						IsBlocked: false,
						DiscordID: "",
					}, nil)
				store.EXPECT().
					SetSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			name: "InternalServerError/DBGetSession",
			buildStubs: func(store *mockmemdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(memdb.Session{}, sql.ErrConnDone)
				store.EXPECT().
					SetSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			name: "InternalServerError/DBSetSession",
			buildStubs: func(store *mockmemdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(session, nil)
				store.EXPECT().
					SetSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(memdb.Session{}, sql.ErrConnDone)
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

			memStore := mockmemdb.NewMockStore(ctrl)
			tc.buildStubs(memStore)

			router := gin.New()
			authMiddleware := middlewares.NewAuthMiddleware(tokenMaker)
			authController := NewAuthController(nil, memStore, config, tokenMaker)
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
