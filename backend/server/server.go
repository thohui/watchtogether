package server

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/thohui/watchtogether/server/routes"
	"github.com/thohui/watchtogether/store"
	"github.com/thohui/watchtogether/youtube"
)

type Server struct {
	app   *fiber.App
	store *store.Store
}

func New() *Server {
	app := fiber.New()
	store := store.New()
	service := youtube.New(os.Getenv("YOUTUBE_API_KEY"))
	routes.Setup(app, store, service)
	return &Server{app: app, store: store}
}

func (server *Server) Start() {
	server.app.Listen(":80")
}

func (server *Server) Stop() {
	server.app.Shutdown()
}
