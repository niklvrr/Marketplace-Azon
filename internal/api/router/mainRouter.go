package router

import (
	"github.com/niklvrr/myMarketplace/internal/handler/cartHandler"
	"github.com/niklvrr/myMarketplace/internal/handler/categoriesHandler"
	"github.com/niklvrr/myMarketplace/internal/handler/orderHandler"
	"github.com/niklvrr/myMarketplace/internal/handler/productHandler"
	"github.com/niklvrr/myMarketplace/internal/handler/userHandler"
	"github.com/niklvrr/myMarketplace/internal/service/cartService"
	"github.com/niklvrr/myMarketplace/internal/service/categoriesService"
	"github.com/niklvrr/myMarketplace/internal/service/orderService"
	"github.com/niklvrr/myMarketplace/internal/service/productService"
	"github.com/niklvrr/myMarketplace/internal/service/userService"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niklvrr/myMarketplace/internal/config"
	"github.com/niklvrr/myMarketplace/internal/repository"
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
	productService := productService.NewProductService(productRepo, rdb)
	userService := userService.NewUserService(userRepo, rdb, jwtManager)
	categoryService := categoriesService.NewCategoriesService(categoryRepo)
	cartService := cartService.NewCartService(cartRepo)
	orderService := orderService.NewOrderService(orderRepo)

	// Handler init
	productHandler := productHandler.NewProductsHandler(productService)
	userHandler := userHandler.NewUserHandler(userService)
	categoryHandler := categoriesHandler.NewCategoryHandler(categoryService)
	cartHandler := cartHandler.NewCartHandler(cartService)
	orderHandler := orderHandler.NewOrderHandler(orderService)

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
