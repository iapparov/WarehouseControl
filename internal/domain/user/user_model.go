package user

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"errors"
	"time"
)

type Role string

const (
	Admin   Role = "admin"
	Manager Role = "manager"
	Viewer  Role = "viewer"
)

type User struct {
	Id        uuid.UUID
	Login     string
	Password  []byte
	Role      Role
	CreatedAt time.Time
}

func NewUser(login, password string, role Role) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if role != Admin && role != Manager && role != Viewer {
		return nil, errors.New("invalid role type")
	}
	if err != nil {
		return nil, err
	}
	return &User{
		Id:        uuid.New(),
		Login:     login,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		Role:      role,
	}, nil
}
