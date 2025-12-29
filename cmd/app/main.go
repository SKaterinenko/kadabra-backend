package main

import (
	"context"
	"errors"
	"fmt"
	"kadabra/internal/config"
	"kadabra/internal/delivery/http/categoryHandler"
	"kadabra/internal/delivery/http/manufacturerHandler"
	"kadabra/internal/delivery/http/productHandler"
	"kadabra/internal/delivery/http/productsTypeHandler"
	"kadabra/internal/delivery/http/subCategoryHandler"
	repository "kadabra/internal/repository/postgres"
	"kadabra/internal/service/categoryService"
	"kadabra/internal/service/manufacturerService"
	"kadabra/internal/service/productService"
	"kadabra/internal/service/productsTypeService"
	"kadabra/internal/service/subCategoryService"
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	// Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("config error: ", err)
	}

	// Database
	ctx := context.Background()
	postgresDB, err := repository.NewPostgres(ctx, cfg)
	if err != nil {
		log.Fatal("db error: ", err)
	}

	// Repository
	categoryRepository := repository.NewCategoryPostgres(postgresDB)
	subCategoryRepository := repository.NewSubCategoryPostgres(postgresDB)
	manufacturerRepository := repository.NewManufacturerPostgres(postgresDB)
	productsTypeRepository := repository.NewProductsTypePostgres(postgresDB)
	productRepository := repository.NewProductPostgres(postgresDB)

	// Service
	category := categoryService.NewService(categoryRepository)
	subCategory := subCategoryService.NewService(subCategoryRepository)
	manufacturer := manufacturerService.NewService(manufacturerRepository)
	productsType := productsTypeService.NewService(productsTypeRepository)
	product := productService.NewService(productRepository)

	// Handlers
	categoryHandler.NewHandler(router, &categoryHandler.HandlerDeps{
		Service: category,
	})
	subCategoryHandler.NewHandler(router, &subCategoryHandler.HandlerDeps{
		Service: subCategory,
	})
	manufacturerHandler.NewHandler(router, &manufacturerHandler.HandlerDeps{
		Service: manufacturer,
	})
	productsTypeHandler.NewHandler(router, &productsTypeHandler.HandlerDeps{
		Service: productsType,
	})
	productHandler.NewHandler(router, &productHandler.HandlerDeps{
		Service: product,
	})

	fmt.Println("Config", cfg)
	fmt.Println("Server is listening on port", cfg.SERVER_PORT)
	server := &http.Server{
		Addr:    cfg.SERVER_PORT,
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server error: %v\n", err)
	}
}
