package routers

import (
	"strings"

	"warehousecontrol/internal/domain/user"
	"warehousecontrol/internal/web/handlers"

	wbgin "github.com/wb-go/wbf/ginext"
)

const (
	CtxUserID = "userId"
	CtxRole   = "role"
	CtxLogin  = "login"
)

func AuthMiddleware(userService handlers.UserIFace) wbgin.HandlerFunc {
	return func(c *wbgin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(401, wbgin.H{"error": "missing token"})
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		if token == "" {
			c.AbortWithStatusJSON(401, wbgin.H{"error": "empty token"})
			return
		}

		payload, err := userService.ValidateTokens(token)
		if err != nil {
			c.AbortWithStatusJSON(401, wbgin.H{"error": "invalid token"})
			return
		}

		c.Set(CtxUserID, payload.UserID)
		c.Set(CtxRole, payload.Role)
		c.Set(CtxLogin, payload.Login)

		c.Next()
	}
}

func RequireRoles(roles ...user.Role) wbgin.HandlerFunc {
	return func(c *wbgin.Context) {

		strRole, exists := c.Get(CtxRole)
		if !exists {
			c.AbortWithStatusJSON(500, wbgin.H{"error": "role not found in context"})
			return
		}

		role, ok := strRole.(user.Role)
		if !ok {
			c.AbortWithStatusJSON(500, wbgin.H{"error": "invalid role type in context"})
			return
		}

		for _, r := range roles {
			if r == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(403, wbgin.H{"error": "forbidden"})
	}
}
