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
	"github.com/redis/go-redis/v9"
)

func NewRouter(db *pgxpool.Pool, rdb *redis.Client, JWTConfig config.JWTConfig) http.Handler {
	// Repository init
	productRepo := repository.NewProductRepo(db)
	userRepo := repository.NewUserRepo(db)
	categoryRepo := repository.NewCategoryRepo(db)
	cartRepo := repository.NewCartRepo(db)
	orderRepo := repository.NewOrderRepo(db)

	// JWTManager init
	jwtManager := jwt.NewJWTManager(JWTConfig.Secret, JWTConfig.Expiration)

	// Service init
	productService := service.NewProductService(productRepo)
	userService := service.NewUserService(userRepo, rdb, jwtManager)
	categoryService := service.NewCategoriesService(categoryRepo)
	cartService := service.NewCartService(cartRepo)
	orderService := service.NewOrderService(orderRepo)

	// Handler init
	productHandler := handler.NewProductsHandler(productService)
	userHandler := handler.NewUserHandler(userService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	cartHandler := handler.NewCartHandler(cartService)
	orderHandler := handler.NewOrderHandler(orderService)

	r := gin.Default()

	api := r.Group("/api")
	v1 := api.Group("/v1")

	registerProductRouter(v1, productHandler, jwtManager, rdb)
	registerUserRouter(v1, userHandler, jwtManager, rdb)
	registerCategoriesRouter(v1, categoryHandler, jwtManager, rdb)
	registerCartRouter(v1, cartHandler, jwtManager, rdb)
	registerOrderRouter(v1, orderHandler, jwtManager, rdb)

	return r
}
