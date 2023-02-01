package service

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"trade-bot/pkg/client/app"
	"trade-bot/pkg/client/models"
)

var (
	ErrSignIn = errors.New("sign in")
	ErrSignUp = errors.New("sign up")
	ErrLogout = errors.New("logout")
)

type AuthService struct {
	client app.ClientActions
}

func NewAuthService(client app.ClientActions) *AuthService {
	return &AuthService{client: client}
}

func (s *AuthService) SignUp(input models.SignUpInput) (models.SignUpResponse, error) {
	req, err := s.client.NewRequest(http.MethodPost, "/auth/sign-up", "", input)
	if err != nil {
		return models.SignUpResponse{}, fmt.Errorf("%s: %w", ErrSignUp, err)
	}

	var output models.SignUpResponse

	resp, err := s.client.Do(req, &output)
	if err != nil {
		return models.SignUpResponse{}, fmt.Errorf("%s: %w", ErrSignUp, err)
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 400) {
		return models.SignUpResponse{}, fmt.Errorf("%s: %s: %s", ErrSignUp, resp.Status, output.Message)
	}

	return output, err
}

func (s *AuthService) SignIn(input models.SignInInput) (models.SignInResponse, error) {
	req, err := s.client.NewRequest(http.MethodPost, "/auth/sign-in", "", input)
	if err != nil {
		return models.SignInResponse{}, fmt.Errorf("%s: %w", ErrSignIn, err)
	}

	var output models.SignInResponse

	resp, err := s.client.Do(req, &output)
	if err != nil {
		return models.SignInResponse{}, fmt.Errorf("%s: %w", ErrSignIn, err)
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 400) {
		return models.SignInResponse{}, fmt.Errorf("%s: %s: %s", ErrSignIn, resp.Status, output.Message)
	}

	return output, err
}

func (s *AuthService) Logout(input models.LogoutInput) (models.LogoutResponse, error) {
	req, err := s.client.NewRequest(http.MethodDelete, "/auth/logout", input.JWTToken, nil)
	if err != nil {
		return models.LogoutResponse{}, fmt.Errorf("%s: %w", ErrLogout, err)
	}

	var output models.LogoutResponse

	resp, err := s.client.Do(req, &output)
	if err != nil {
		return models.LogoutResponse{}, fmt.Errorf("%s: %w", ErrLogout, err)
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 400) {
		return models.LogoutResponse{}, fmt.Errorf("%s: %s: %s", ErrLogout, resp.Status, output.Message)
	}

	return output, err
}
