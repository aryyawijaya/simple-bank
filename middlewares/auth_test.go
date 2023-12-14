package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/aryyawijaya/simple-bank/middlewares"
	"github.com/aryyawijaya/simple-bank/modules/auth/token"
	testhelper "github.com/aryyawijaya/simple-bank/util/test-helper"
	"github.com/aryyawijaya/simple-bank/util/wrapper"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

var w = wrapper.NewWrapper()
var middleware = middlewares.NewMiddleware(w)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				testhelper.CustomAuthorization(
					t,
					request,
					tokenMaker,
					middlewares.AuthorizationTypeBearer,
					"username",
					time.Minute,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "No authorization header",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Invalid authorization header format",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				testhelper.CustomAuthorization(
					t,
					request,
					tokenMaker,
					"",
					"username",
					time.Minute,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Unsupported authorization type",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				testhelper.CustomAuthorization(
					t,
					request,
					tokenMaker,
					"unsupported",
					"username",
					time.Minute,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Expired token",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				testhelper.CustomAuthorization(
					t,
					request,
					tokenMaker,
					middlewares.AuthorizationTypeBearer,
					"username",
					-1*time.Minute,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// setup server test
			server := testhelper.NewTestServer(t, nil)

			authPath := "/auth"
			server.Router.GET(
				authPath,
				middleware.AuthMiddleware(server.TokenMaker),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			// request
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.TokenMaker)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
