package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Aqui se incializara los route de la api
func SetupRoutes(app *fiber.App, db *pgxpool.Pool) {

	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("API Running")
	})
}
