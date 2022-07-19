package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thohui/watchtogether/room"
	"github.com/thohui/watchtogether/store"
)

type body struct {
	VideoId string `json:"video_id"`
}

type createRoomResponse struct {
	Id string `json:"id"`
}

func createRoomHandler(store *store.Store) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		body := new(body)
		if err := c.BodyParser(body); err != nil || body.VideoId == "" {
			return c.Status(400).SendString("Invalid request")
		}
		// TODO: validate video id using the youtube api
		room := room.New(body.VideoId)
		store.Add(room)
		return c.Status(200).JSON(createRoomResponse{Id: room.Id})
	}
}
