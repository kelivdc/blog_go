package main

import (
	"blog/database"
	"blog/middlewares"
	"blog/routes"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		middlewares.SendError("Error loading .env file")
	}

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	database.InitDatabase()
	routes.InitRoutes(app)
	log.Fatal(app.Listen(":8080"))
}
