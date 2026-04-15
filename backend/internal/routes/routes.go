package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Aqui se definiran los endpoints de la API
func SetupRoutes(app *fiber.App, db *pgxpool.Pool) {

	app.Post("/upload", func(c fiber.Ctx) error {
		return c.SendString("Subiendo archivos")
	})

	app.Post("/download/:token", func(c fiber.Ctx) error {
		return c.SendString("Descargando archivo")
	})

	app.Get("/file/:token", func(c fiber.Ctx) error {
		return c.SendString("Obteninedo archivo")
	})

	app.Delete("/file/:token", func(c fiber.Ctx) error {
		return c.SendString("Eliminando archivo")
	})

	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("API Running")
	})
}
