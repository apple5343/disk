package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	type respBody struct {
		Message string `json:"message"`
	}
	tests := []struct {
		name    string
		handler echo.HandlerFunc
		resp    respBody
		code    int
		wantErr bool
	}{
		{
			name: "success",
			handler: func(_ echo.Context) error {
				return nil
			},
			code: http.StatusOK,
		},
		{
			name: "response commited",
			handler: func(c echo.Context) error {
				return c.JSON(http.StatusCreated, map[string]string{"message": "ok"})
			},
			code: http.StatusCreated,
		},
		{
			name: "error unauthorized",
			handler: func(_ echo.Context) error {
				return errorx.NewError("unauthorized", errorx.Unauthorized)
			},
			resp: respBody{
				Message: "unauthorized",
			},
			code:    http.StatusUnauthorized,
			wantErr: true,
		},
		{
			name: "error internal",
			handler: func(_ echo.Context) error {
				return errorx.NewError("some internal error", errorx.Internal)
			},
			resp: respBody{
				Message: "Internal server error",
			},
			code:    http.StatusInternalServerError,
			wantErr: true,
		},
		{
			name: "error unknown",
			handler: func(_ echo.Context) error {
				return errors.New("unknown")
			},
			resp: respBody{
				Message: "Internal server error",
			},
			code:    http.StatusInternalServerError,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			c := echo.New().NewContext(req, rec)

			ErrorMiddleware()(tt.handler)(c)
			require.Equal(t, tt.code, rec.Code)
			if tt.wantErr {
				var resp respBody
				require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
				require.Equal(t, tt.resp.Message, resp.Message)
			}
		})
	}
}
