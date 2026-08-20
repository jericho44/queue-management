package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"queue-management-tenant/backend/internal/config"
	ws "queue-management-tenant/backend/internal/websocket"
	"queue-management-tenant/backend/pkg/database"
	"queue-management-tenant/backend/pkg/jwt"
	appLogger "queue-management-tenant/backend/pkg/logger"
	appRedis "queue-management-tenant/backend/pkg/redis"
	"queue-management-tenant/backend/pkg/response"

	// Domain Modules
	auditCtrl "queue-management-tenant/backend/internal/modules/audit/controller"
	auditRepo "queue-management-tenant/backend/internal/modules/audit/repository"
	auditRoutes "queue-management-tenant/backend/internal/modules/audit/routes"

	authCtrl "queue-management-tenant/backend/internal/modules/auth/controller"
	authRepo "queue-management-tenant/backend/internal/modules/auth/repository"
	authRoutes "queue-management-tenant/backend/internal/modules/auth/routes"
	authSvc "queue-management-tenant/backend/internal/modules/auth/service"

	branchCtrl "queue-management-tenant/backend/internal/modules/branch/controller"
	branchRepo "queue-management-tenant/backend/internal/modules/branch/repository"
	branchRoutes "queue-management-tenant/backend/internal/modules/branch/routes"
	branchSvc "queue-management-tenant/backend/internal/modules/branch/service"

	counterCtrl "queue-management-tenant/backend/internal/modules/counter/controller"
	counterRepo "queue-management-tenant/backend/internal/modules/counter/repository"
	counterRoutes "queue-management-tenant/backend/internal/modules/counter/routes"
	counterSvc "queue-management-tenant/backend/internal/modules/counter/service"

	orgRepo "queue-management-tenant/backend/internal/modules/organization/repository"

	queueCtrl "queue-management-tenant/backend/internal/modules/queue/controller"
	queueRepo "queue-management-tenant/backend/internal/modules/queue/repository"
	queueRoutes "queue-management-tenant/backend/internal/modules/queue/routes"
	queueSvc "queue-management-tenant/backend/internal/modules/queue/service"

	reportCtrl "queue-management-tenant/backend/internal/modules/report/controller"
	reportRepo "queue-management-tenant/backend/internal/modules/report/repository"
	reportRoutes "queue-management-tenant/backend/internal/modules/report/routes"
	reportSvc "queue-management-tenant/backend/internal/modules/report/service"

	serviceCtrl "queue-management-tenant/backend/internal/modules/service/controller"
	serviceRepo "queue-management-tenant/backend/internal/modules/service/repository"
	serviceRoutes "queue-management-tenant/backend/internal/modules/service/routes"
	serviceSvc "queue-management-tenant/backend/internal/modules/service/service"
)

func main() {
	cfg := config.LoadConfig()
	appLog := appLogger.NewLogger()

	// Connect PostgreSQL
	db, err := database.NewPostgresConnection(database.PostgresConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	})
	if err != nil {
		appLog.Error("Database connection warning: %v", err)
	} else {
		appLog.Info("Connected to PostgreSQL database: %s", cfg.DBName)
	}

	// Connect Redis
	redisClient, err := appRedis.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword, 0)
	if err != nil {
		appLog.Error("Redis connection warning: %v (falling back to in-memory pubsub)", err)
	} else {
		appLog.Info("Connected to Redis at %s", cfg.RedisAddr)
	}

	// JWT Service
	jwtSvc := jwt.NewJWTService(cfg.JWTSecret, 24*1000000000*60*60, 7*24*1000000000*60*60)

	// WebSocket Hub
	wsHub := ws.NewWSHub(redisClient, appLog)
	go wsHub.Run(context.Background())

	// Module Repositories
	oRepo := orgRepo.NewOrganizationRepository(db)
	uRepo := authRepo.NewUserRepository(db)
	bRepo := branchRepo.NewBranchRepository(db)
	sRepo := serviceRepo.NewServiceRepository(db)
	cRepo := counterRepo.NewCounterRepository(db)
	qRepo := queueRepo.NewQueueRepository(db)
	rRepo := reportRepo.NewReportRepository(db)
	aRepo := auditRepo.NewAuditRepository(db)

	// Module Services
	aService := authSvc.NewAuthService(db, oRepo, uRepo, jwtSvc)
	bService := branchSvc.NewBranchService(bRepo, oRepo)
	sService := serviceSvc.NewServiceManagementService(sRepo)
	cService := counterSvc.NewCounterService(cRepo, oRepo)
	qService := queueSvc.NewQueueService(qRepo, sRepo, cRepo, wsHub)
	rService := reportSvc.NewReportService(rRepo)

	// Module Controllers
	aController := authCtrl.NewAuthController(aService)
	bController := branchCtrl.NewBranchController(bService)
	sController := serviceCtrl.NewServiceController(sService)
	cController := counterCtrl.NewCounterController(cService)
	qController := queueCtrl.NewQueueController(qService)
	rController := reportCtrl.NewReportController(rService)
	auditController := auditCtrl.NewAuditController(aRepo)

	// Fiber App Initialization
	app := fiber.New(fiber.Config{
		AppName:      "Queue Management System SaaS API",
		ServerHeader: "Go-Fiber",
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CorsOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Request-ID",
	}))

	// Health Check Endpoints
	app.Get("/health", func(c *fiber.Ctx) error {
		return response.Success(c, fiber.StatusOK, "System operational", fiber.Map{
			"status": "UP",
			"db":     db != nil,
			"redis":  redisClient != nil,
		})
	})
	app.Get("/ready", func(c *fiber.Ctx) error {
		return response.Success(c, fiber.StatusOK, "Ready", nil)
	})

	// WebSocket Endpoint
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws", websocket.New(wsHub.HandleWS))

	// API V1 Group
	api := app.Group("/api/v1")
	api.Get("/", func(c *fiber.Ctx) error {
		return response.Success(c, fiber.StatusOK, "Queue Management System SaaS API v1 Operational", nil)
	})

	// Register Module Routes
	authRoutes.RegisterAuthRoutes(api, aController, jwtSvc)
	branchRoutes.RegisterBranchRoutes(api, bController, jwtSvc)
	serviceRoutes.RegisterServiceRoutes(api, sController, jwtSvc)
	counterRoutes.RegisterCounterRoutes(api, cController, jwtSvc)
	queueRoutes.RegisterQueueRoutes(api, qController, jwtSvc)
	reportRoutes.RegisterReportRoutes(api, rController, jwtSvc)
	auditRoutes.RegisterAuditRoutes(api, auditController, jwtSvc)

	port := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Starting Queue Management System API server on port %s", port)
	if err := app.Listen(port); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
