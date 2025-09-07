package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/niklvrr/myMarketplace/internal/errs"
	"github.com/niklvrr/myMarketplace/pkg/jwt"
	"github.com/redis/go-redis/v9"
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

func JWTRegister(jwtManager *jwt.JWTManager, cache *redis.Client) gin.HandlerFunc {
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

		blacklistKey := "blacklist_user:" + strconv.Itoa(int(claims.UserId))
		exist, err := cache.Exists(c.Request.Context(), blacklistKey).Result()
		if err != nil {
			errs.RespondError(c, http.StatusUnauthorized, "unauthorized", err.Error())
			c.Abort()
			return
		}

		if exist > 0 {
			errs.RespondError(c, http.StatusForbidden, "forbidden", "token has been revoked")
			c.Abort()
			return
		}

		blockKey := "blocked_user:" + strconv.Itoa(int(claims.UserId))
		exist, err = cache.Exists(c.Request.Context(), blockKey).Result()
		if err != nil {
			errs.RespondError(c, http.StatusUnauthorized, "unauthorized", err.Error())
			c.Abort()
			return
		}

		if exist > 0 {
			errs.RespondError(c, http.StatusForbidden, "forbidden", "user has been blocked")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserId)
		c.Set("role", claims.Role)
		c.Next()
	}
}
