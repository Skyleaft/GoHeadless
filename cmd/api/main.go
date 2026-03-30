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
	"log"
	"os"

		"GoHeadless/docs"
	"GoHeadless/internal/collection"
	"GoHeadless/internal/content"
	"GoHeadless/internal/platform"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
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
	defer mongoDB.Client.Disconnect(nil)

	// 3. Initialize Repositories
	collRepo := collection.NewRepository(mongoDB.Database)
	recordRepo := content.NewRepository(mongoDB.Database)

	// 4. Initialize Services
	collService := collection.NewService(collRepo)
	contentService := content.NewService(recordRepo, collRepo)

	// 5. Initialize Handlers
	collHandler := collection.NewHandler(collService)
	contentHandler := content.NewHandler(contentService)

	// 6. Setup Fiber Application
	app := fiber.New(fiber.Config{
		AppName: "GoHeadless CMS v1",
	})

	// Middleware
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(recover.New())

	// 7. Define Routes
	app.Get("/docs/*", scalar.New())

	api := app.Group("/api/v1")
	collHandler.Routes(api)
	contentHandler.Routes(api)

	// Start Server
	log.Printf("GoHeadless CMS is running on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
