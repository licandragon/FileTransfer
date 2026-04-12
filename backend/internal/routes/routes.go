package routes

import "github.com/gofiber/fiber/v3"

//Aqui se declararan lo
func SetupRoutes(app *fiber.App) {

	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("API Running")
	})
}
