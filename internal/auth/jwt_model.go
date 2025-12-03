package auth

import (
	"warehousecontrol/internal/domain/user"
)

type JWTResponse struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresIn  int64
	RefreshExpiresIn int64
	TokenType        string
}

type JWTPayload struct {
	UserID string
	Role   user.Role
	Login  string
}
