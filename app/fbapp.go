package fbapp

import "github.com/gofiber/fiber/v2"

var App *fiber.App

func init() {
	App = fiber.New()
}
