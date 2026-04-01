package main

// @title GoHeadless CMS API
// @version 1.0
// @description This is a dynamic Headless CMS API built with Go Fiber and MongoDB.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api/v1

import (
	"context"
	"errors"
	"log"
	"os"

	"GoHeadless/docs"
	"GoHeadless/internal/apierr"
	"GoHeadless/internal/auth"
	"GoHeadless/internal/collection"
	"GoHeadless/internal/content"
	"GoHeadless/internal/middleware"
	"GoHeadless/internal/platform"
	"GoHeadless/internal/setup"
	"GoHeadless/internal/upload"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/joho/godotenv"
	"github.com/yokeTH/gofiber-scalar/scalar/v3"
)

func main() {
	// 0. Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// 1. Initial configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	docs.SwaggerInfo.Host = "localhost:" + port
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "goheadless_cms"
	}

	// 2. Setup Platform (MongoDB)
	mongoDB := platform.NewMongoDB(mongoURI, dbName)
	defer mongoDB.Client.Disconnect(context.Background())

	// 3. Initialize Repositories
	collRepo := collection.NewRepository(mongoDB.Database)
	recordRepo := content.NewRepository(mongoDB.Database)
	authRepo := auth.NewRepository(mongoDB.Database)

	if err := collRepo.EnsurePhysicalCollections(context.Background()); err != nil {
		log.Printf("ensure physical MongoDB collections: %v", err)
	}

	// 4. Initialize Services
	collService := collection.NewService(collRepo)
	contentService := content.NewService(recordRepo, collRepo)
	uploadService := upload.NewService("./uploads")
	authService := auth.NewService(authRepo, collRepo, recordRepo)
	setupService := setup.NewService(authRepo)

	// 5. Initialize Handlers
	collHandler := collection.NewHandler(collService)
	contentHandler := content.NewHandler(contentService)
	uploadHandler := upload.NewHandler(uploadService)
	authHandler := auth.NewHandler(authService)
	setupHandler := setup.NewHandler(setupService)

	// 6. Initialize Middleware
	rbac := middleware.NewRBACMiddleware(authService, collService)

	// 6. Setup Fiber Application
	app := fiber.New(fiber.Config{
		AppName: "GoHeadless CMS v1",
		ErrorHandler: func(c fiber.Ctx, err error) error {
			var ae *apierr.AppError
			if ok := errors.As(err, &ae); ok {
				if ae.Internal != nil {
					log.Printf("[ERROR] %s %s → %s: %v", c.Method(), c.Path(), ae.Message, ae.Internal)
				}
				return c.Status(ae.Code).JSON(fiber.Map{"error": ae.Message})
			}

			// fiber.Error (e.g. 404 from router)
			if fe, ok := err.(*fiber.Error); ok {
				log.Printf("[FIBER] %s %s → %d %s", c.Method(), c.Path(), fe.Code, fe.Message)
				return nil
				//return c.Status(fe.Code).JSON(fiber.Map{"error": fe.Message})
			}

			// Unexpected error — log full detail, return generic message
			log.Printf("[ERROR] %s %s → unhandled: %v", c.Method(), c.Path(), err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		},
	})

	// Middleware
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(recover.New())

	// 7. Define Routes
	app.Get("/docs/*", scalar.New())

	// Static files for uploaded images
	app.Use("/uploads", static.New("./uploads"))

	api := app.Group("/api/v1")

	// Public Setup & Auth
	setupHandler.Routes(api)
	authHandler.PublicRoutes(api)

	// Content routes: AuthorizeCollection handles public access check first,
	// then authentication if needed. NOT under protected group to allow public access.
	contentGroup := api.Group("/content/:slug", rbac.AuthorizeContentCollection)
	contentHandler.Routes(contentGroup)

	// Protected Routes (Require Authentication)
	protected := api.Group("", rbac.Authenticate)

	// Collection & Content with RBAC
	// Collections still require authentication
	collectionGroup := protected.Group("/collections", rbac.AuthorizeCollection)
	collHandler.Routes(collectionGroup)

	uploadHandler.Routes(protected.Group("/upload"))

	// Admin-only Routes (Superadmin required)
	adminGroup := protected.Group("/admin", rbac.RequireSuperadmin)
	authHandler.AdminRoutes(adminGroup)

	// Start Server
	log.Printf("GoHeadless CMS is running on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
