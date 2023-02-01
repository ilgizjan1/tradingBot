package utils

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
)

const (
	accessUUIDTokenClaim = "access_UUID"
	authorizedTokenClaim = "authorized"
	userIDTokenClaim     = "user_id"
	expiresTokenClaim    = "exp"
)

const jwtHeaderAlgo = "alg"
const authorizationHeader = "Authorization"
const jwtAccessSigningKey = "JWT_ACCESS_SIGNING_KEY"

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrVerifyToken             = errors.New("verify token")
	ErrExtractTokenMetadata    = errors.New("extract token metadata")
	ErrCantAssignToMapClaims   = errors.New("can't assign to map claims")
	ErrInvalidAccessUUID       = errors.New("invalid access uuid")
	ErrInvalidUserID           = errors.New("invalid user id")
	ErrInvalidToken            = errors.New("invalid token")
	ErrEmptyAuthHeader         = errors.New("empty auth header")
	ErrInvalidAuthHeader       = errors.New("invalid auth header")
	ErrEmptyBearerToken        = errors.New("empty bearer token")
)

type TokenDetails struct {
	AccessToken string
	AccessUUID  string
	AtExpires   int64
}

type AccessDetails struct {
	AccessUUID string
	UserID     int64
}

func GetBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get(authorizationHeader)

	if authHeader == "" {
		return "", ErrEmptyAuthHeader
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return "", ErrInvalidAuthHeader
	}

	if headerParts[0] != "Bearer" {
		return "", ErrInvalidAuthHeader
	}

	if len(headerParts[1]) == 0 {
		return "", ErrEmptyBearerToken
	}

	return headerParts[1], nil
}

func ExtractTokenMetadata(token string) (AccessDetails, error) {
	verifiedToken, err := VerifyToken(token)
	if err != nil {
		return AccessDetails{}, fmt.Errorf("%s: %w", ErrExtractTokenMetadata, err)
	}

	claims, ok := verifiedToken.Claims.(jwt.MapClaims)
	if !ok {
		return AccessDetails{}, fmt.Errorf("%s: %w", ErrExtractTokenMetadata, ErrCantAssignToMapClaims)
	}

	accessUUID, ok := claims[accessUUIDTokenClaim].(string)
	if !ok {
		return AccessDetails{}, fmt.Errorf("%s: %w", ErrExtractTokenMetadata, ErrInvalidAccessUUID)
	}
	userID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims[userIDTokenClaim]), 10, 64)
	if err != nil {
		return AccessDetails{}, fmt.Errorf("%s: %w", ErrExtractTokenMetadata, ErrInvalidUserID)
	}

	return AccessDetails{
		AccessUUID: accessUUID,
		UserID:     int64(userID),
	}, nil
}

func VerifyToken(token string) (*jwt.Token, error) {
	verified, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: %v", ErrUnexpectedSigningMethod, token.Header[jwtHeaderAlgo])
		}
		return []byte(os.Getenv(jwtAccessSigningKey)), nil
	})
	if err != nil {
		return nil, err
	}

	if ok := verified.Valid; !ok {
		return nil, fmt.Errorf("%s: %w", ErrVerifyToken, ErrInvalidToken)
	}
	return verified, nil
}

func GenerateJWTToken(userID int, d time.Duration) (TokenDetails, error) {
	td := TokenDetails{}

	td.AtExpires = time.Now().Add(d).Unix()
	aUUID, err := uuid.NewV4()
	if err != nil {
		return TokenDetails{}, err
	}
	td.AccessUUID = aUUID.String()

	atClaims := jwt.MapClaims{}
	atClaims[authorizedTokenClaim] = true
	atClaims[accessUUIDTokenClaim] = td.AccessUUID
	atClaims[userIDTokenClaim] = userID
	atClaims[expiresTokenClaim] = td.AtExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv(jwtAccessSigningKey)))

	return td, err
}
