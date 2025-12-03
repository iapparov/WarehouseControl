package handlers

import (
	"net/http"

	"warehousecontrol/internal/auth"
	"warehousecontrol/internal/domain/user"
	"warehousecontrol/internal/web/dto"

	wbgin "github.com/wb-go/wbf/ginext"
)

type UserHandler struct {
	Service UserIFace
}

type UserIFace interface {
	Login(Login, Password string) (*auth.JWTResponse, error)
	Registration(Login, Password, Role string) (*user.User, error)
	RefreshTokens(tokenStr string) (*auth.JWTResponse, error)
	ValidateTokens(tokenStr string) (*auth.JWTPayload, error)
}

func NewUserHandler(service UserIFace) *UserHandler {
	return &UserHandler{
		Service: service,
	}
}

// RegisterUser
// @Summary Register a new user
// @Description Create a new user account with role
// @Tags users
// @Accept json
// @Produce json
// @Param body body dto.UserRegistrationRequest true "User registration info"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/register [post]
func (h *UserHandler) RegisterUser(ctx *wbgin.Context) {
	var req dto.UserRegistrationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, wbgin.H{"error": err.Error()})
		return
	}
	user, err := h.Service.Registration(req.Login, req.Password, req.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, wbgin.H{"error": err.Error()})
		return
	}
	res := dto.UserResponse{
		ID:    user.Id.String(),
		Login: user.Login,
		Role:  string(user.Role),
	}
	ctx.JSON(http.StatusOK, res)
}

// LoginUser
// @Summary Login user
// @Description Authenticate user and return JWT tokens
// @Tags users
// @Accept json
// @Produce json
// @Param body body dto.UserLoginRequest true "User login info"
// @Success 200 {object} dto.JWTResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/login [post]
func (h *UserHandler) LoginUser(ctx *wbgin.Context) {
	var req dto.UserLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, wbgin.H{"error": err.Error()})
		return
	}
	jwtResp, err := h.Service.Login(req.Login, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, wbgin.H{"error": err.Error()})
		return
	}
	res := dto.JWTResponse{
		AccessToken:  jwtResp.AccessToken,
		RefreshToken: jwtResp.RefreshToken,
	}
	ctx.JSON(http.StatusOK, res)
}

// RefreshToken
// @Summary Refresh JWT token
// @Description Refresh access and refresh tokens using existing refresh token
// @Tags users
// @Accept json
// @Produce json
// @Param body body dto.TokenRefreshRequest true "Refresh token request"
// @Success 200 {object} dto.JWTResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/refresh-token [post]
func (h *UserHandler) RefreshToken(ctx *wbgin.Context) {
	var req dto.TokenRefreshRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, wbgin.H{"error": err.Error()})
		return
	}
	jwtResp, err := h.Service.RefreshTokens(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, wbgin.H{"error": err.Error()})
		return
	}
	res := dto.JWTResponse{
		AccessToken:  jwtResp.AccessToken,
		RefreshToken: jwtResp.RefreshToken,
	}
	ctx.JSON(http.StatusOK, res)
}
