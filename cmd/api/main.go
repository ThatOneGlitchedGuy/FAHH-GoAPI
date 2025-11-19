package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang-app/internal/api/handlers"
	"golang-app/internal/api/middleware"
	"golang-app/internal/config"
	"golang-app/internal/database"
	"golang-app/internal/repository"
	"golang-app/internal/service"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg, err := config.GetSettings()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if cfg.Debug {
		log.Printf("Loaded configuration: %+v\n", cfg)
	}

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	service.InitCache()

	service.InitScheduler()
	service.StartJobs()

	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()
	router.Use(middleware.MetricsMiddleware())
	pprof.Register(router)

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Advanced Go API!"})
	})
	
	apiV1 := router.Group(cfg.APIV1Prefix)
	
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	statsRepo := repository.NewStatsRepository(db)
	
	userService := service.NewUserService(userRepo, cfg)
	authService := service.NewAuthService(userRepo, cfg)
	productService := service.NewProductService(productRepo, orderRepo, reviewRepo, userRepo, cfg)
	orderService := service.NewOrderService(orderRepo, productRepo, cfg)
	messageService := service.NewMessageService(messageRepo, orderRepo, cfg)
	statsService := service.NewStatsService(statsRepo)
	
	authHandler := handlers.NewAuthHandler(authService, userService)
	userHandler := handlers.NewUserHandler(userService)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)
	messageHandler := handlers.NewMessageHandler(messageService)
	adminHandler := handlers.NewAdminHandler(statsService)
	
	authHandler.RegisterRoutes(apiV1.Group("/auth"))
	userHandler.RegisterRoutes(apiV1.Group("/users"))
	productHandler.RegisterRoutes(apiV1.Group("/products"))
	orderHandler.RegisterRoutes(apiV1.Group("/orders"))
	messageHandler.RegisterRoutes(apiV1.Group("/messages"))
	adminHandler.RegisterRoutes(apiV1)
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
