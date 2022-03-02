package main

import (
	"log"

	"github.com/colcrunch/avecalc_backend/database"
	"github.com/colcrunch/avecalc_backend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func welcome(c *fiber.Ctx) error {
	return c.SendString("Welcome!")
}

func main() {
	app := fiber.New()
	app.Use(cors.New())

	database.ConnectDB()

	app.Get("/", welcome)
	routes.UserRoutes(app)
	routes.AuthRoutes(app)
	routes.ContractRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
