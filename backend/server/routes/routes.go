package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/thohui/watchtogether/store"
	"github.com/thohui/watchtogether/youtube"
)

func Setup(app *fiber.App, store *store.Store, youtube *youtube.YoutubeService) {
	app.Use(cors.New())
	app.Use("/ws", websocketMiddleware)
	app.Get("/ws/:id", websocketHandler(store))
	app.Post("/room/create", createRoomHandler(store, youtube))
	app.Post("/room/:id", getRoomHandler(store))
}
