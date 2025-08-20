package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niklvrr/myMarketplace/internal/config"
	"github.com/niklvrr/myMarketplace/internal/handler"
	"github.com/niklvrr/myMarketplace/internal/repository"
	"github.com/niklvrr/myMarketplace/internal/service"
	"github.com/niklvrr/myMarketplace/pkg/jwt"
)

func NewRouter(db *pgxpool.Pool, JWTConfig config.JWTConfig) http.Handler {
	// Repository init
	productRepo := repository.NewProductRepo(db)
	userRepo := repository.NewUserRepo(db)
	//categoryRepo := repository.NewCategoryRepo(db)
	//cartRepo := repository.NewCartRepo(db)
	//orderRepo := repository.NewOrderRepo(db)
	//adminRepo := repository.NewAdminRepo(db)

	// JWTManager init
	jwtManager := jwt.NewJWTManager(JWTConfig.Secret, JWTConfig.Expiration)

	// Service init
	productService := service.NewProductService(productRepo)
	userService := service.NewUserService(userRepo, jwtManager)

	// Handler init
	productHandler := handler.NewProductsHandler(productService)
	userHandler := handler.NewUserHandler(userService)

	r := gin.Default()

	api := r.Group("/api")
	v1 := api.Group("/v1")

	registerProductRouter(v1, productHandler, jwtManager)
	registerUserRouter(v1, userHandler, jwtManager)

	return r
}
