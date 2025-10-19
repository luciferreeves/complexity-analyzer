package main

import (
	"complexity-analyzer/handlers"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize template engine
	engine := html.New("./templates", ".html")

	app := fiber.New(fiber.Config{
		Views:     engine,
		BodyLimit: 10 * 1024 * 1024,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Static files
	app.Static("/css", "./static/css")
	app.Static("/js", "./static/js")

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Algorithm Complexity Analyzer",
		})
	})

	// API routes
	api := app.Group("/api")
	api.Post("/analyze", handlers.AnalyzeCode)

	log.Println("Server starting on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}
