package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/thohui/watchtogether/store"
)

func Setup(app *fiber.App, store *store.Store) {
	app.Use(cors.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Use("/ws", websocketMiddleware)
	app.Get("/ws/:id", websocketHandler(store))
	app.Post("/room/create", createRoomHandler(store))
}
