package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niklvrr/myMarketplace/internal/handlers"
	"github.com/niklvrr/myMarketplace/internal/repository"
	"github.com/niklvrr/myMarketplace/internal/service"
)

func NewRouter(db *pgxpool.Pool) http.Handler {
	// Repository init
	productRepo := repository.NewProductRepo(db)
	//userRepo := repository.NewUserRepo(db)
	//categoryRepo := repository.NewCategoryRepo(db)
	//cartRepo := repository.NewCartRepo(db)
	//orderRepo := repository.NewOrderRepo(db)
	//adminRepo := repository.NewAdminRepo(db)

	// Service init
	productService := service.NewProductService(productRepo)

	// Handler init
	productHandler := handlers.NewProductsHandler(productService)

	r := gin.Default()

	api := r.Group("/api")
	v1 := api.Group("/v1")

	registerProductRouter(v1, productHandler)

	return r
}
