package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thohui/watchtogether/store"
)

func Setup(app *fiber.App, store *store.Store) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Use("/ws", websocketMiddleware)
	app.Get("/ws/:id", websocketRoute(store))
}
