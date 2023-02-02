package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thohui/watchtogether/room"
	"github.com/thohui/watchtogether/store"
	"github.com/thohui/watchtogether/youtube"
)

type createRoomResponse struct {
	Id string `json:"id"`
}

func createRoomHandler(store *store.Store, youtube *youtube.YoutubeService) func(c *fiber.Ctx) error {
	type body struct {
		VideoId string `json:"video_id"`
	}
	return func(c *fiber.Ctx) error {
		body := new(body)
		if err := c.BodyParser(body); err != nil || body.VideoId == "" {
			return c.Status(400).SendString("Invalid request")
		}
		video, err := youtube.GetVideo(body.VideoId)
		if err != nil {
			return c.Status(400).SendString("Invalid video id")
		}
		room := room.New(video)
		store.Add(room)
		return c.Status(200).JSON(createRoomResponse{Id: room.Id})
	}
}

func getRoomHandler(store *store.Store) func(c *fiber.Ctx) error {
	type body struct {
		Id string `json:"id"`
	}
	return func(c *fiber.Ctx) error {
		body := new(body)
		if err := c.BodyParser(body); err != nil || body.Id == "" {
			return c.Status(400).SendString("Invalid request")
		}
		_, err := store.Get(body.Id)
		if err != nil {
			return c.Status(404).SendString(err.Error())
		}
		return c.Status(200).JSON(body)
	}
}
