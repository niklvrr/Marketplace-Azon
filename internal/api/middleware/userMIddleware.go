package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/niklvrr/myMarketplace/internal/errs"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{})
	for _, role := range roles {
		allowed[role] = struct{}{}
	}
	return func(c *gin.Context) {
		role, ok := c.Get("role")
		if !ok {
			errs.RespondError(c, http.StatusForbidden, "forbidden", "user don`t have role access")
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
