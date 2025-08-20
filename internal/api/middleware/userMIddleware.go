package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/niklvrr/myMarketplace/internal/errs"
	"github.com/niklvrr/myMarketplace/pkg/jwt"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{})
	for _, role := range roles {
		allowed[role] = struct{}{}
	}
	return func(c *gin.Context) {
		role, ok := c.Get("role")
		if !ok {
			errs.RespondError(c, http.StatusForbidden, "forbidden", "no info about role in context")
			c.Abort()
			return
		}

		if _, ok = allowed[role.(string)]; !ok {
			errs.RespondError(c, http.StatusForbidden, "forbidden", "user don`t have role access")
			c.Abort()
			return
		}

		c.Next()
	}
}

func JWTRegister(jwtManager *jwt.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			errs.RespondError(c, http.StatusForbidden, "forbidden", "user is not authorized")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			errs.RespondError(c, http.StatusForbidden, "forbidden", "invalid token")
			c.Abort()
			return
		}

		claims, err := jwtManager.ParseToken(tokenString)
		if err != nil {
			errs.RespondError(c, http.StatusForbidden, "forbidden", err.Error())
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserId)
		c.Set("role", claims.Role)
		c.Next()
	}
}
