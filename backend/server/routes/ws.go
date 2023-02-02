package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/thohui/watchtogether/store"
)

func websocketMiddleware(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired

}

func websocketHandler(store *store.Store) func(c *fiber.Ctx) error {
	return websocket.New(func(c *websocket.Conn) {
		id := c.Params("id")
		room, err := store.Get(id)
		if err != nil {
			c.Close()
			return
		}
		room.Handle(c.Conn)
	})
}
