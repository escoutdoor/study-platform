package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/escoutdoor/study-platform/internal/apperror"
	"github.com/escoutdoor/study-platform/pkg/errwrap"
	"github.com/golang-jwt/jwt/v5"
)

type Role int32

const (
	RoleStudent Role = iota + 1
	RoleTeacher
)

type TokenProvider struct {
	accessTokenSecretKey  string
	refreshTokenSecretKey string
	accessTokenTTL        time.Duration
	refreshTokenTTL       time.Duration
}

func NewTokenProvider(
	accessTokenSecretKey string,
	refreshTokenSecretKey string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *TokenProvider {
	return &TokenProvider{
		accessTokenSecretKey:  accessTokenSecretKey,
		refreshTokenSecretKey: refreshTokenSecretKey,
		accessTokenTTL:        accessTokenTTL,
		refreshTokenTTL:       refreshTokenTTL,
	}
}

type AccessTokenClaims struct {
	jwt.RegisteredClaims

	UserID int    `json:"userId"`
	Roles  []Role `json:"roles"`
}

type RefreshTokenClaims struct {
	jwt.RegisteredClaims

	UserID int `json:"userId"`
}

func (p *TokenProvider) GenerateAccessToken(userID int, roles []Role) (string, error) {
	claims := AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.accessTokenTTL)),
		},
		UserID: userID,
		Roles:  roles,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(p.accessTokenSecretKey))
	if err != nil {
		return "", errwrap.Wrap("new jwt token with claims", err)
	}

	return token, nil
}

func (p *TokenProvider) GenerateRefreshToken(userID int) (string, error) {
	claims := RefreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.refreshTokenTTL)),
		},
		UserID: userID,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(p.refreshTokenSecretKey))
	if err != nil {
		return "", errwrap.Wrap("new jwt token with claims", err)
	}

	return token, nil
}

func (p *TokenProvider) ValidateAccessToken(accessToken string) (AccessTokenClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(accessToken, &AccessTokenClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(p.accessTokenSecretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return AccessTokenClaims{}, apperror.ErrJwtTokenExpired
		}

		return AccessTokenClaims{}, apperror.ErrInvalidJwtToken
	}

	if !jwtToken.Valid {
		return AccessTokenClaims{}, apperror.ErrInvalidJwtToken
	}

	claims, ok := jwtToken.Claims.(*AccessTokenClaims)
	if !ok {
		return AccessTokenClaims{}, apperror.ErrInvalidJwtToken
	}

	return *claims, nil
}

func (p *TokenProvider) ValidateRefreshToken(refreshToken string) (RefreshTokenClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(refreshToken, &RefreshTokenClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(p.refreshTokenSecretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return RefreshTokenClaims{}, apperror.ErrJwtTokenExpired
		}

		return RefreshTokenClaims{}, apperror.ErrInvalidJwtToken
	}

	if !jwtToken.Valid {
		return RefreshTokenClaims{}, apperror.ErrInvalidJwtToken
	}

	claims, ok := jwtToken.Claims.(*RefreshTokenClaims)
	if !ok {
		return RefreshTokenClaims{}, apperror.ErrInvalidJwtToken
	}

	return *claims, nil
}
