package user

import (
	"testing"
)

func TestNewUser_ValidRoles(t *testing.T) {
	for _, role := range []Role{Admin, Manager, Viewer} {
		u, err := NewUser("john", "secret123", role)
		if err != nil {
			t.Fatalf("expected no error for role %s, got %v", role, err)
		}
		if u.Login != "john" || u.Role != role {
			t.Fatalf("unexpected user fields: %+v", u)
		}
		if len(u.Password) == 0 {
			t.Fatalf("expected hashed password")
		}
	}
}

func TestNewUser_InvalidRole(t *testing.T) {
	_, err := NewUser("john", "secret123", Role("bad"))
	if err == nil {
		t.Fatalf("expected error for invalid role")
	}
}
