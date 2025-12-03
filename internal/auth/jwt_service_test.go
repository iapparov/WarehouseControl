package auth_test

import (
	"testing"
	"time"

	"warehousecontrol/internal/auth"
	"warehousecontrol/internal/config"
	"warehousecontrol/internal/domain/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func newTestConfig() *config.AppConfig {
	return &config.AppConfig{
		JwtConfig: config.JwtConfig{
			JwtAccessSecret:    "access-secret",
			JwtRefreshSecret:   "refresh-secret",
			JwtExpAccessToken:  1, // 1 минута
			JwtExpRefreshToken: 1, // 1 час
		},
	}
}

func newTestJWT() *auth.JWTService {
	return auth.NewJWTService(newTestConfig())
}

func newTestUser() *user.User {
	return &user.User{
		Id:    uuid.New(),
		Login: "testlogin",
		Role:  user.Role("admin"),
	}
}

func TestGenerateTokens_Success(t *testing.T) {
	s := newTestJWT()
	u := newTestUser()

	resp, err := s.GenerateTokens(u)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.AccessToken == "" || resp.RefreshToken == "" {
		t.Fatal("tokens must not be empty")
	}
}

func TestValidateTokens_Success(t *testing.T) {
	s := newTestJWT()
	u := newTestUser()

	tokens, _ := s.GenerateTokens(u)

	payload, err := s.ValidateTokens(tokens.AccessToken)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if payload.UserID != u.Id.String() {
		t.Fatalf("invalid payload user ID")
	}
	if payload.Login != u.Login {
		t.Fatalf("invalid payload login")
	}
	if payload.Role != u.Role {
		t.Fatalf("invalid payload role")
	}
}

func TestValidateTokens_InvalidToken(t *testing.T) {
	s := newTestJWT()

	_, err := s.ValidateTokens("this_is_invalid_token")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestValidateTokens_NoUUIDField(t *testing.T) {
	s := newTestJWT()

	claims := jwt.MapClaims{
		"login": "test",
		"role":  "admin",
		"exp":   time.Now().Add(time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte("access-secret"))

	_, err := s.ValidateTokens(tokenStr)
	if err == nil {
		t.Fatal("expected error because uuid field is missing")
	}
}

func TestValidateTokens_NoLoginOrRole(t *testing.T) {
	s := newTestJWT()

	claims := jwt.MapClaims{
		"uuid": uuid.New().String(),
		"exp":  time.Now().Add(time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte("access-secret"))

	_, err := s.ValidateTokens(tokenStr)
	if err == nil {
		t.Fatal("expected error because login/role missing")
	}
}

func TestRefreshTokens_Success(t *testing.T) {
	s := newTestJWT()
	u := newTestUser()

	tokens, _ := s.GenerateTokens(u)

	newTokens, err := s.RefreshTokens(tokens.RefreshToken)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if newTokens.AccessToken == "" || newTokens.RefreshToken == "" {
		t.Fatal("new tokens must not be empty")
	}
}

func TestRefreshTokens_Invalid(t *testing.T) {
	s := newTestJWT()

	_, err := s.RefreshTokens("invalid_refresh_token")
	if err == nil {
		t.Fatal("expected error for invalid refresh token")
	}
}

func TestRefreshTokens_NoUUIDField(t *testing.T) {
	s := newTestJWT()

	claims := jwt.MapClaims{
		"login": "test",
		"role":  "admin",
		"exp":   time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte("refresh-secret"))

	_, err := s.RefreshTokens(tokenStr)
	if err == nil {
		t.Fatal("expected error because uuid missing")
	}
}

func TestRefreshTokens_NoLoginOrRole(t *testing.T) {
	s := newTestJWT()

	claims := jwt.MapClaims{
		"uuid": uuid.New().String(),
		"exp":  time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte("refresh-secret"))

	_, err := s.RefreshTokens(tokenStr)
	if err == nil {
		t.Fatal("expected error because login/role missing")
	}
}

func TestValidateAccessToken_Expired(t *testing.T) {
	cfg := &config.AppConfig{
		JwtConfig: config.JwtConfig{
			JwtAccessSecret:   "access-secret",
			JwtRefreshSecret:  "refresh-secret",
			JwtExpAccessToken: -1, // уже истёк
		},
	}

	s := auth.NewJWTService(cfg)
	u := newTestUser()

	expired, _ := s.GenerateTokens(u)

	_, err := s.ValidateTokens(expired.AccessToken)
	if err == nil {
		t.Fatal("expected error: access token expired")
	}
}

func TestValidateRefreshToken_Expired(t *testing.T) {
	cfg := &config.AppConfig{
		JwtConfig: config.JwtConfig{
			JwtAccessSecret:    "access-secret",
			JwtRefreshSecret:   "refresh-secret",
			JwtExpRefreshToken: -1, // час назад истёк
		},
	}

	s := auth.NewJWTService(cfg)
	u := newTestUser()

	tokens, _ := s.GenerateTokens(u)

	_, err := s.RefreshTokens(tokens.RefreshToken)
	if err == nil {
		t.Fatal("expected error: refresh token expired")
	}
}
