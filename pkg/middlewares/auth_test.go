package middlewares

import (
	"fmt"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/modules/token"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, r *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				accessToken, _, err := tokenMaker.CreateToken("1234", time.Minute)
				require.NoError(t, err)

				authHeader := fmt.Sprintf("%s %s", AuthorizationTypeBearer, accessToken)
				r.Header.Set(AuthorizationHeaderKey, authHeader)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:      "NoAuthorization",
			setupAuth: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			name: "UnsupportedAuthorizationType",
			setupAuth: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				accessToken, _, err := tokenMaker.CreateToken("1234", time.Minute)
				require.NoError(t, err)

				authHeader := fmt.Sprintf("%s %s", "unsupported", accessToken)
				r.Header.Set(AuthorizationHeaderKey, authHeader)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setupAuth: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				accessToken, _, err := tokenMaker.CreateToken("1234", time.Minute)
				require.NoError(t, err)

				authHeader := fmt.Sprintf("%s %s", "", accessToken)
				r.Header.Set(AuthorizationHeaderKey, authHeader)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				accessToken, _, err := tokenMaker.CreateToken("1234", -time.Minute)
				require.NoError(t, err)

				authHeader := fmt.Sprintf("%s %s", AuthorizationTypeBearer, accessToken)
				r.Header.Set(AuthorizationHeaderKey, authHeader)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()

			tokenMaker, _ := token.NewPasetoMaker(utils.RandomString(32))
			authMiddleware := NewAuthMiddleware(tokenMaker)
			router.GET("/auth", authMiddleware, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{})
			})

			req, err := http.NewRequest(http.MethodGet, "/auth", nil)
			require.NoError(t, err)

			tc.setupAuth(t, req, tokenMaker)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			tc.checkResponse(t, w)
		})
	}
}
