package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/BoggerByte/Sentinel-backend.git/pkg/db/mock"
	db "github.com/BoggerByte/Sentinel-backend.git/pkg/db/sqlc"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/middlewares"
	token2 "github.com/BoggerByte/Sentinel-backend.git/pkg/modules/token"
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

func generateRandomUser() db.User {
	return db.User{
		ID:            int64(utils.RandomInt(1, 1000)),
		DiscordID:     utils.RandomSnowflakeID().String(),
		Username:      gofakeit.Username(),
		Discriminator: fmt.Sprintf("%04d", utils.RandomInt(1, 9999)),
		Verified:      gofakeit.Bool(),
		Email:         gofakeit.Email(),
		Avatar:        gofakeit.ImageURL(400, 400),
		Banner:        gofakeit.ImageURL(400, 400),
		AccentColor:   int64(utils.RandomInt(0, 1<<24)),
		CreatedAt:     gofakeit.Date(),
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)

	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user, gotUser)
}

func TestUserController_Get(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	user := generateRandomUser()

	testCases := []struct {
		name          string
		userDiscordID string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:          "OK",
			userDiscordID: user.DiscordID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.DiscordID)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				requireBodyMatchUser(t, w.Body, user)
			},
		},
		{
			name:          "NotFound",
			userDiscordID: user.DiscordID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.DiscordID)).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			name:          "InternalServerError",
			userDiscordID: user.DiscordID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.DiscordID)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
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
			// build server
			tokenMaker, _ := token2.NewPasetoMaker(utils.RandomString(32))
			authMiddleware := middlewares.NewAuthMiddleware(tokenMaker)
			userController := NewUserController(store)
			router := gin.New()
			router.GET("/api/v1/users/me", authMiddleware, userController.GetUser)

			url := "/api/v1/users/me"
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			accessToken, _, err := tokenMaker.CreateToken(tc.userDiscordID, time.Minute)
			require.NoError(t, err)

			authHeader := fmt.Sprintf("%s %s", middlewares.AuthorizationTypeBearer, accessToken)
			req.Header.Set(middlewares.AuthorizationHeaderKey, authHeader)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			tc.checkResponse(t, w)
		})
	}
}
