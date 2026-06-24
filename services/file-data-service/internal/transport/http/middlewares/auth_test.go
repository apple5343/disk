package middlewares

import (
	"data/internal/config"
	ctxutil "data/internal/utils/ctx"
	"data/pkg/jwt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	jwtConfig := &config.JWTConfig{
		Secret: "secret",
	}
	ttl := time.Hour
	userID := gofakeit.UUID()
	tests := []struct {
		name    string
		request func() *http.Request
		code    int
		wantErr bool
	}{
		{
			name: "success",
			request: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				token, err := jwt.GenerateToken(jwt.UserInfo{ID: userID}, []byte(jwtConfig.Secret), ttl)
				require.NoError(t, err)
				req.Header.Set("Authorization", "Bearer "+token)
				return req
			},
			code: http.StatusOK,
		},
		{
			name: "invalid token",
			request: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "Bearer invalid_token")
				return req
			},
			code:    http.StatusUnauthorized,
			wantErr: true,
		},
		{
			name: "invalid header",
			request: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "Invalid token")
				return req
			},
			code:    http.StatusUnauthorized,
			wantErr: true,
		},
		{
			name: "no token",
			request: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				return req
			},
			code:    http.StatusUnauthorized,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.request()
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)
			ErrorMiddleware()(AuthMiddleware(jwtConfig)(func(c echo.Context) error {
				if !tt.wantErr {
					id := ctxutil.UserIDFromContext(c.Request().Context())
					require.Equal(t, userID, id)
				}
				return nil
			}))(c)

			require.Equal(t, tt.code, rec.Code)
		})
	}
}
