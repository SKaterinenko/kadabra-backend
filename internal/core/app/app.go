package app

import (
	"context"
	"kadabra/internal/core/config"
	categories_http "kadabra/internal/features/categories/delivery/http"
	categories_postgres "kadabra/internal/features/categories/repository/postgres"
	categories_service "kadabra/internal/features/categories/service"
	manufacturers_http "kadabra/internal/features/manufacturers/delivery/http"
	manufacturers_postgres "kadabra/internal/features/manufacturers/repository/postgres"
	manufacturers_service "kadabra/internal/features/manufacturers/service"
	products_http "kadabra/internal/features/products/delivery/http"
	products_postgres "kadabra/internal/features/products/repository/postgres"
	products_service "kadabra/internal/features/products/service"
	products_type_http "kadabra/internal/features/products_type/delivery/http"
	products_type_postgres "kadabra/internal/features/products_type/repository/postgres"
	products_type_service "kadabra/internal/features/products_type/service"
	sub_categories_http "kadabra/internal/features/sub_categories/delivery/http"
	sub_categories_postgres "kadabra/internal/features/sub_categories/repository/postgres"
	sub_categories_service "kadabra/internal/features/sub_categories/service"
	users_http "kadabra/internal/features/users/handler/http"
	users_postgres "kadabra/internal/features/users/repository"
	users_service "kadabra/internal/features/users/service"
	"kadabra/pkg/middleware"

	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func App() (context.Context, http.Handler, *config.Config, func()) {
	router := http.NewServeMux()

	// Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("config error: ", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	// Database
	postgresDB, err := config.NewPostgres(ctx, cfg)
	if err != nil {
		stop()
		log.Fatal("db error: ", err)
	}

	// S3
	s3Client, err := config.NewS3Client(cfg)
	if err != nil {
		log.Fatal("s3 error: ", err)
	}

	// Repository
	categoryRepository := categories_postgres.NewCategoryPostgres(postgresDB)
	subCategoryRepository := sub_categories_postgres.NewSubCategoryPostgres(postgresDB)
	manufacturerRepository := manufacturers_postgres.NewManufacturerPostgres(postgresDB)
	productsTypeRepository := products_type_postgres.NewProductsTypePostgres(postgresDB)
	productsRepository := products_postgres.NewProductPostgres(postgresDB)
	usersRepository := users_postgres.NewUsersPostgres(postgresDB)

	// Service
	category := categories_service.NewService(categoryRepository)
	subCategory := sub_categories_service.NewService(subCategoryRepository)
	manufacturer := manufacturers_service.NewService(manufacturerRepository)
	productsType := products_type_service.NewService(productsTypeRepository)
	products := products_service.NewService(productsRepository, s3Client)
	users := users_service.NewService(usersRepository, s3Client, cfg)

	// Handlers
	categories_http.NewHandler(router, &categories_http.HandlerDeps{
		Service: category,
	})
	sub_categories_http.NewHandler(router, &sub_categories_http.HandlerDeps{
		Service: subCategory,
	})
	manufacturers_http.NewHandler(router, &manufacturers_http.HandlerDeps{
		Service: manufacturer,
	})
	products_type_http.NewHandler(router, &products_type_http.HandlerDeps{
		Service: productsType,
	})
	products_http.NewHandler(router, &products_http.HandlerDeps{
		Service: products,
	})
	users_http.NewHandler(router, &users_http.HandlerDeps{
		Service: users,
		Cfg:     cfg,
	})

	// Middlewares
	stack := middleware.Chain(
		middleware.CORS,
	)

	cleanup := func() {
		stop()
		postgresDB.Close()
	}

	return ctx, stack(router), cfg, cleanup
}
