package user_test

import (
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"warehousecontrol/internal/app/user"
	"warehousecontrol/internal/auth"
	"warehousecontrol/internal/config"
	domain "warehousecontrol/internal/domain/user"
)

type fakeRepo struct {
	users map[string]*domain.User
	err   error
}

func (f *fakeRepo) GetUser(login string) (*domain.User, error) {
	if f.err != nil {
		return nil, f.err
	}
	u, ok := f.users[login]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (f *fakeRepo) SaveUser(u *domain.User) error {
	if f.err != nil {
		return f.err
	}
	if f.users == nil {
		f.users = map[string]*domain.User{}
	}
	f.users[u.Login] = u
	return nil
}

type fakeJwt struct{}

func (f *fakeJwt) GenerateTokens(u *domain.User) (*auth.JWTResponse, error) {
	return &auth.JWTResponse{AccessToken: "access", RefreshToken: "refresh"}, nil
}
func (f *fakeJwt) ValidateTokens(tokenStr string) (*auth.JWTPayload, error) {
	return &auth.JWTPayload{UserID: "uid", Login: "login"}, nil
}
func (f *fakeJwt) RefreshTokens(refreshToken string) (*auth.JWTResponse, error) {
	return &auth.JWTResponse{AccessToken: "access", RefreshToken: "refresh"}, nil
}

func testCfg() *config.AppConfig {
	return &config.AppConfig{
		UserConfig: config.UserConfig{
			MinLength:         3,
			MaxLength:         10,
			AllowedCharacters: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-",
		},
		PasswordConfig: config.PasswordConfig{
			MinLength:    6,
			MaxLength:    12,
			RequireUpper: true,
			RequireLower: true,
			RequireDigit: true,
		},
	}
}

func TestLogin_Success(t *testing.T) {
	pass := "Password1"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	u := &domain.User{Login: "user", Password: hashed}

	repo := &fakeRepo{users: map[string]*domain.User{"user": u}}
	jwt := &fakeJwt{}
	svc := user.NewUserService(repo, jwt, testCfg())

	tokens, err := svc.Login("user", pass)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		t.Fatalf("expected tokens")
	}
}

func TestLogin_Errors(t *testing.T) {
	repo := &fakeRepo{users: map[string]*domain.User{}}
	jwt := &fakeJwt{}
	svc := user.NewUserService(repo, jwt, testCfg())

	if _, err := svc.Login("", "pass"); err == nil {
		t.Fatal("expected error for empty login")
	}

	if _, err := svc.Login("user", ""); err == nil {
		t.Fatal("expected error for empty password")
	}

	if _, err := svc.Login("unknown", "pass"); err == nil {
		t.Fatal("expected error for unknown user")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte("Password1"), bcrypt.DefaultCost)
	repo.users["user"] = &domain.User{Login: "user", Password: hashed}
	if _, err := svc.Login("user", "wrongpass"); err == nil {
		t.Fatal("expected error for wrong password")
	}
}

func TestRegistration_Success(t *testing.T) {
	repo := &fakeRepo{}
	jwt := &fakeJwt{}
	svc := user.NewUserService(repo, jwt, testCfg())

	u, err := svc.Registration("valid", "Password1", "viewer")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Login != "valid" {
		t.Fatalf("unexpected login")
	}
}

func TestRegistration_Errors(t *testing.T) {
	repo := &fakeRepo{users: map[string]*domain.User{"exist": {Login: "exist"}}}
	jwt := &fakeJwt{}
	svc := user.NewUserService(repo, jwt, testCfg())

	if _, err := svc.Registration("ab", "Password1", "viewer"); err == nil {
		t.Fatal("expected error for login too short")
	}

	if _, err := svc.Registration("bad$", "Password1", "viewer"); err == nil {
		t.Fatal("expected error for invalid chars")
	}

	if _, err := svc.Registration("exist", "Password1", "viewer"); err == nil {
		t.Fatal("expected error for existing user")
	}

	if _, err := svc.Registration("newuser", "pass", "viewer"); err == nil {
		t.Fatal("expected error for invalid password")
	}
}

func TestIsValidLogin(t *testing.T) {
	repo := &fakeRepo{}
	svc := user.NewUserService(repo, &fakeJwt{}, testCfg())

	_, err := svc.Registration("ab", "Password1", "viewer")
	if err == nil {
		t.Fatal("expected error for short login")
	}

	_, err = svc.Registration("bad$", "Password1", "viewer")
	if err == nil {
		t.Fatal("expected error for invalid chars")
	}

	u, err := svc.Registration("good123", "Password1", "viewer")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Login != "good123" {
		t.Fatalf("expected login to be good123, got %s", u.Login)
	}
}

func TestIsValidPassword(t *testing.T) {
	repo := &fakeRepo{}
	svc := user.NewUserService(repo, &fakeJwt{}, testCfg())

	tests := []struct {
		pwd string
		ok  bool
	}{
		{"pass", false},              // too short
		{"verylongpassword1", false}, // too long
		{"lowercase1", false},        // no upper
		{"UPPERCASE1", false},        // no lower
		{"NoDigit", false},           // no digit
		{"Valid1", true},
	}

	for _, tt := range tests {

		_, err := svc.Registration("user123", tt.pwd, "viewer")
		if tt.ok && err != nil && err.Error() != "user with this login already exists" {
			t.Fatalf("expected success for pwd %s, got %v", tt.pwd, err)
		}
		if !tt.ok && err == nil {
			t.Fatalf("expected error for pwd %s", tt.pwd)
		}
	}
}

func TestRefreshAndValidateTokens(t *testing.T) {
	svc := user.NewUserService(nil, &fakeJwt{}, testCfg())

	r, err := svc.RefreshTokens("refresh")
	if err != nil || r.AccessToken == "" {
		t.Fatal("expected valid refresh token response")
	}

	v, err := svc.ValidateTokens("token")
	if err != nil || v.UserID == "" {
		t.Fatal("expected valid validate token response")
	}
}
