package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niklvrr/myMarketplace/internal/models"
	"net/http"
)

// User alias for user model
type User = models.User

func signUp(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var r User
		if err := c.ShouldBindJSON(&r); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

	}
}
