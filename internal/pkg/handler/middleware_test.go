package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"trade-bot/internal/pkg/service"
	mockService "trade-bot/internal/pkg/service/mocks"
)

func TestHandler_userIdentity(t *testing.T) {
	type mockBehaviour func(s *mockService.MockAuthorization, token string)

	tests := []struct {
		name                          string
		headerName                    string
		headerValue                   string
		token                         string
		mockBehaviourOnGetUserID      mockBehaviour
		mockBehaviourOnGetUserAPIKeys mockBehaviour
		expectedStatusCode            int
		expectedRequestBody           string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehaviourOnGetUserID: func(s *mockService.MockAuthorization, token string) {
				s.EXPECT().GetUserIDByJWT(token).Return(1, nil)
			},
			mockBehaviourOnGetUserAPIKeys: func(s *mockService.MockAuthorization, token string) {
				s.EXPECT().GetUserAPIKeys(1).Return("public", "private", nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{1, public, private}`,
		},
		{
			name:                          "Invalid header name",
			headerName:                    "",
			headerValue:                   "Bearer token",
			mockBehaviourOnGetUserID:      func(s *mockService.MockAuthorization, token string) {},
			mockBehaviourOnGetUserAPIKeys: func(s *mockService.MockAuthorization, token string) {},
			expectedStatusCode:            http.StatusUnauthorized,
			expectedRequestBody:           `{"message":"user identity: empty auth header"}`,
		},
		{
			name:                          "Invalid header value",
			headerName:                    "Authorization",
			headerValue:                   "Bearerrrrrrrr token",
			mockBehaviourOnGetUserID:      func(s *mockService.MockAuthorization, token string) {},
			mockBehaviourOnGetUserAPIKeys: func(s *mockService.MockAuthorization, token string) {},
			expectedStatusCode:            http.StatusUnauthorized,
			expectedRequestBody:           `{"message":"user identity: invalid auth header"}`,
		},
		{
			name:                          "Empty token",
			headerName:                    "Authorization",
			headerValue:                   "Bearer ",
			mockBehaviourOnGetUserID:      func(s *mockService.MockAuthorization, token string) {},
			mockBehaviourOnGetUserAPIKeys: func(s *mockService.MockAuthorization, token string) {},
			expectedStatusCode:            http.StatusUnauthorized,
			expectedRequestBody:           `{"message":"user identity: empty bearer token"}`,
		},
		{
			name:        "Service error on GetUserIDByJWT",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehaviourOnGetUserID: func(s *mockService.MockAuthorization, token string) {
				s.EXPECT().GetUserIDByJWT(token).Return(0, errors.New("bad token"))
			},
			mockBehaviourOnGetUserAPIKeys: func(s *mockService.MockAuthorization, token string) {},
			expectedStatusCode:            http.StatusUnauthorized,
			expectedRequestBody:           `{"message":"user identity: bad token"}`,
		},
		{
			name:        "Service error on GetUserAPIKeys",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehaviourOnGetUserID: func(s *mockService.MockAuthorization, token string) {
				s.EXPECT().GetUserIDByJWT(token).Return(1, nil)
			},
			mockBehaviourOnGetUserAPIKeys: func(s *mockService.MockAuthorization, token string) {
				s.EXPECT().GetUserAPIKeys(1).Return("", "", errors.New("bad token"))
			},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedRequestBody: `{"message":"user identity: bad token"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockService.NewMockAuthorization(c)
			test.mockBehaviourOnGetUserID(repo, test.token)
			test.mockBehaviourOnGetUserAPIKeys(repo, test.token)

			services := &service.Service{Authorization: repo}
			handler := Handler{services, nil, nil}

			r := gin.New()
			r.GET("/identity", handler.userIdentity, func(c *gin.Context) {
				id, _ := c.Get(userIDCtx)
				publicKey, _ := c.Get(userPublicAPIKeyCtx)
				privateKey, _ := c.Get(userPrivateAPIKeyCtx)
				c.String(http.StatusOK, "{%d, %s, %s}", id, publicKey, privateKey)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/identity", nil)
			req.Header.Set(test.headerName, test.headerValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_getUserID(t *testing.T) {
	setContext := func(value interface{}) *gin.Context {
		ctx := &gin.Context{}
		ctx.Set(userIDCtx, value)
		return ctx
	}

	tests := []struct {
		name    string
		ctx     *gin.Context
		wantID  int
		wantErr bool
	}{
		{
			name:   "OK",
			ctx:    setContext(1),
			wantID: 1,
		},
		{
			name:    "Empty context",
			ctx:     &gin.Context{},
			wantErr: true,
		},
		{
			name:    "Invalid value in context",
			ctx:     setContext("invalid id"),
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			id, err := getUserID(test.ctx)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.wantID, id)
			}
		})
	}
}

func TestHandler_getUserAPIKeys(t *testing.T) {
	setContext := func(publicKey, privateKey string) *gin.Context {
		ctx := &gin.Context{}
		if publicKey != "" {
			ctx.Set(userPublicAPIKeyCtx, publicKey)
		}
		if privateKey != "" {
			ctx.Set(userPrivateAPIKeyCtx, privateKey)
		}
		return ctx
	}

	tests := []struct {
		name        string
		ctx         *gin.Context
		wantPublic  string
		wantPrivate string
		wantErr     bool
	}{
		{
			name:        "OK",
			ctx:         setContext("public", "private"),
			wantPublic:  "public",
			wantPrivate: "private",
		},
		{
			name:    "Empty context",
			ctx:     &gin.Context{},
			wantErr: true,
		},
		{
			name:    "Empty public",
			ctx:     setContext("", "private"),
			wantErr: true,
		},
		{
			name:    "Empty private",
			ctx:     setContext("public", ""),
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			public, private, err := getUserAPIKeys(test.ctx)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.wantPublic, public)
				assert.Equal(t, test.wantPrivate, private)
			}
		})
	}
}
