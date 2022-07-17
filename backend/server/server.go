package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thohui/watchtogether/server/routes"
	"github.com/thohui/watchtogether/store"
)

type Server struct {
	app   *fiber.App
	store *store.Store
}

func New() *Server {
	app := fiber.New()
	store := store.New()
	routes.Setup(app, store)
	return &Server{app: app, store: store}
}

func (server *Server) Start() {
	server.app.Listen(":80")
}

func (server *Server) Stop() {
	server.app.Shutdown()
}
