package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"trade-bot/internal/pkg/models"
	"trade-bot/internal/pkg/service"
	mockService "trade-bot/internal/pkg/service/mocks"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehaviour func(s *mockService.MockAuthorization, user models.User)

	tests := []struct {
		name                string
		inputBody           string
		inputUser           models.User
		mockBehaviour       mockBehaviour
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			inputBody: `{
				"name":"name",
				"username":"username",
				"password":"qwerty",
				"public_api_key":"key",
				"private_api_key":"key"
			}`,
			inputUser: models.User{
				Name:          "name",
				Username:      "username",
				Password:      "qwerty",
				PublicAPIKey:  "key",
				PrivateAPIKey: "key",
			},
			mockBehaviour: func(s *mockService.MockAuthorization, user models.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:                "Wrong Input",
			inputBody:           `{"username":"username"}`,
			inputUser:           models.User{},
			mockBehaviour:       func(s *mockService.MockAuthorization, user models.User) {},
			expectedStatusCode:  400,
			expectedRequestBody: fmt.Sprintf(`{"message":"%s"}`, ErrInvalidInputBody),
		},
		{
			name: "Service error",
			inputBody: `{
				"name":"name",
				"username":"username",
				"password":"qwerty",
				"public_api_key":"key",
				"private_api_key":"key"
			}`,
			inputUser: models.User{
				Name:          "name",
				Username:      "username",
				Password:      "qwerty",
				PublicAPIKey:  "key",
				PrivateAPIKey: "key",
			},
			mockBehaviour: func(s *mockService.MockAuthorization, user models.User) {
				s.EXPECT().CreateUser(user).Return(0, errors.New("something went wrong"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockService.NewMockAuthorization(c)
			test.mockBehaviour(repo, test.inputUser)

			services := &service.Service{Authorization: repo}
			handler := Handler{services, nil, nil}

			// test server
			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			// test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/sign-up",
				bytes.NewBufferString(test.inputBody))

			// make request
			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	type mockBehaviour func(s *mockService.MockAuthorization, username, password string)

	tests := []struct {
		name                string
		inputBody           string
		username            string
		password            string
		mockBehaviour       mockBehaviour
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username":"username", "password":"qwerty"}`,
			username:  "username",
			password:  "qwerty",
			mockBehaviour: func(s *mockService.MockAuthorization, username, password string) {
				s.EXPECT().GenerateJWT(username, password).Return("token", nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"access_token":"token"}`,
		},
		{
			name:                "Wrong Input",
			inputBody:           `{"username":"username"}`,
			mockBehaviour:       func(s *mockService.MockAuthorization, username, password string) {},
			expectedStatusCode:  400,
			expectedRequestBody: fmt.Sprintf(`{"message":"%s"}`, ErrInvalidInputBody),
		},
		{
			name:      "Service error",
			inputBody: `{"username":"username", "password":"qwerty"}`,
			username:  "username",
			password:  "qwerty",
			mockBehaviour: func(s *mockService.MockAuthorization, username, password string) {
				s.EXPECT().GenerateJWT(username, password).Return("", errors.New("something went wrong"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockService.NewMockAuthorization(c)
			test.mockBehaviour(repo, test.username, test.password)

			services := &service.Service{Authorization: repo}
			handler := Handler{services, nil, nil}

			r := gin.New()
			r.POST("/sign-in", handler.signIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/sign-in",
				bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_logout(t *testing.T) {
	type mockBehaviour func(s *mockService.MockAuthorization, token string)

	tests := []struct {
		name                string
		headerName          string
		headerValue         string
		token               string
		mockBehaviour       mockBehaviour
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehaviour: func(s *mockService.MockAuthorization, token string) {
				s.EXPECT().LogoutUser(token).Return(nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"message":"successfully logged out"}`,
		},
		{
			name:                "Invalid header name",
			headerName:          "",
			headerValue:         "Bearer token",
			token:               "token",
			mockBehaviour:       func(s *mockService.MockAuthorization, token string) {},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedRequestBody: `{"message":"empty auth header"}`,
		},
		{
			name:                "Invalid header value",
			headerName:          "Authorization",
			headerValue:         "Bearrrrrer token",
			token:               "token",
			mockBehaviour:       func(s *mockService.MockAuthorization, token string) {},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedRequestBody: `{"message":"invalid auth header"}`,
		},
		{
			name:                "Empty token",
			headerName:          "Authorization",
			headerValue:         "Bearer ",
			token:               "token",
			mockBehaviour:       func(s *mockService.MockAuthorization, token string) {},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedRequestBody: `{"message":"empty bearer token"}`,
		},
		{
			name:        "Parse Error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehaviour: func(s *mockService.MockAuthorization, token string) {
				s.EXPECT().LogoutUser(token).Return(errors.New("invalid token"))
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: `{"message":"invalid token"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockService.NewMockAuthorization(c)
			test.mockBehaviour(repo, test.token)

			services := &service.Service{Authorization: repo}
			handler := Handler{services, nil, nil}

			r := gin.New()
			r.DELETE("/logout", handler.logout)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/logout", nil)
			req.Header.Set(test.headerName, test.headerValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedRequestBody, w.Body.String())
		})
	}
}
