package auth

import (
	"errors"
	"time"

	"warehousecontrol/internal/config"
	"warehousecontrol/internal/domain/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	accessSecret       string
	refreshSecret      string
	jwtExpAccessToken  int // в минутах
	jwtExpRefreshToken int // в часах
}

// Конструктор
func NewJWTService(cfg *config.AppConfig) *JWTService {
	return &JWTService{
		accessSecret:       cfg.JwtConfig.JwtAccessSecret,
		refreshSecret:      cfg.JwtConfig.JwtRefreshSecret,
		jwtExpAccessToken:  cfg.JwtConfig.JwtExpAccessToken,
		jwtExpRefreshToken: cfg.JwtConfig.JwtExpRefreshToken,
	}
}

// Генерация пары токенов
func (s *JWTService) GenerateTokens(u *user.User) (*JWTResponse, error) {
	access, err := s.generateAccessToken(u)
	if err != nil {
		return nil, err
	}
	refresh, err := s.generateRefreshToken(u)
	if err != nil {
		return nil, err
	}
	return &JWTResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

// Проверка токена (возвращаем полезную нагрузку)
func (s *JWTService) ValidateTokens(tokenStr string) (*JWTPayload, error) {
	claims, err := s.validateAccessToken(tokenStr)
	if err != nil {
		return nil, err
	}

	uuidStr, ok := claims["uuid"].(string)
	if !ok {
		return nil, errors.New("invalid token payload")
	}
	role, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("invalid token payload")
	}
	login, ok := claims["login"].(string)
	if !ok {
		return nil, errors.New("invalid token payload")
	}
	return &JWTPayload{
		UserID: uuidStr,
		Role:   user.Role(role),
		Login:  login,
	}, nil
}

// Обновление токенов по refresh
func (s *JWTService) RefreshTokens(refreshToken string) (*JWTResponse, error) {
	claims, err := s.validateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	uuidStr, ok := claims["uuid"].(string)
	if !ok {
		return nil, errors.New("invalid refresh token payload")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("invalid refresh token payload")
	}

	login, ok := claims["login"].(string)
	if !ok {
		return nil, errors.New("invalid refresh token payload")
	}

	// Создаём "фейкового" пользователя только для генерации токенов
	u := &user.User{
		Id:    uuid.MustParse(uuidStr),
		Role:  user.Role(role),
		Login: login,
	}

	// Генерируем новую пару токенов
	return s.GenerateTokens(u)
}

//// Вспомогательные приватные методы

func (s *JWTService) generateAccessToken(u *user.User) (string, error) {
	claims := jwt.MapClaims{
		"uuid":  u.Id.String(),
		"role":  u.Role,
		"login": u.Login,
		"exp":   time.Now().Add(time.Minute * time.Duration(s.jwtExpAccessToken)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.accessSecret))
}

func (s *JWTService) generateRefreshToken(u *user.User) (string, error) {
	claims := jwt.MapClaims{
		"uuid":  u.Id.String(),
		"role":  u.Role,
		"login": u.Login,
		"exp":   time.Now().Add(time.Hour * time.Duration(s.jwtExpRefreshToken)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.refreshSecret))
}

func (s *JWTService) validateAccessToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.accessSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid access token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	if exp, ok := claims["exp"].(float64); !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}

func (s *JWTService) validateRefreshToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.refreshSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	if exp, ok := claims["exp"].(float64); !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
		return nil, errors.New("refresh token has expired")
	}

	return claims, nil
}
