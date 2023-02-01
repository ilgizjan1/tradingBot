package service

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"trade-bot/internal/pkg/models"
	"trade-bot/internal/pkg/repository"
	"trade-bot/pkg/utils"
)

var (
	ErrCreateUser         = errors.New("create user")
	ErrGenerateJWT        = errors.New("generate jwt")
	ErrGetUserIDByJWT     = errors.New("get user id by jwt")
	ErrLogoutUser         = errors.New("logout user")
	ErrGetUserAPIKeys     = errors.New("get user api keys")
	ErrMismatchedPassword = errors.New("mismatched password")
)

type AuthService struct {
	repo    repository.Authorization
	jwtRepo repository.JWT
}

func NewAuthService(repo repository.Authorization, jwtRepo repository.JWT) *AuthService {
	return &AuthService{repo: repo, jwtRepo: jwtRepo}
}

func (s *AuthService) CreateUser(user models.User) (int, error) {
	err := user.GeneratePasswordHash(user.Password)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", ErrCreateUser, err)
	}
	userID, err := s.repo.CreateUser(user)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", ErrCreateUser, err)
	}
	return userID, nil
}

func (s *AuthService) GenerateJWT(username string, password string) (string, error) {
	user, err := s.repo.GetUser(username)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrGenerateJWT, err)
	}
	if ok := user.ComparePassword(password); !ok {
		return "", fmt.Errorf("%s: %w", ErrGenerateJWT, ErrMismatchedPassword)
	}

	td, err := utils.GenerateJWTToken(user.ID, time.Hour*12)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrGenerateJWT, err)
	}
	token, err := s.jwtRepo.CreateJWT(user.ID, td)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrGenerateJWT, err)
	}
	return token, nil
}

func (s *AuthService) GetUserIDByJWT(token string) (int, error) {
	ad, err := utils.ExtractTokenMetadata(token)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", ErrGetUserIDByJWT, err)
	}
	userID, err := s.jwtRepo.GetJWTUserID(ad)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", ErrGetUserIDByJWT, err)
	}
	return userID, nil
}

func (s *AuthService) LogoutUser(token string) error {
	ad, err := utils.ExtractTokenMetadata(token)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrLogoutUser, err)
	}
	if err := s.jwtRepo.DeleteJWT(ad); err != nil {
		return fmt.Errorf("%s: %w", ErrLogoutUser, err)
	}
	return nil
}

func (s *AuthService) GetUserAPIKeys(userID int) (string, string, error) {
	public, private, err := s.repo.GetUserAPIKeys(userID)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", ErrGetUserAPIKeys, err)
	}
	return public, private, nil
}
