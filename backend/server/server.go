package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thohui/watchtogether/server/routes"
)

type Server struct {
	app *fiber.App
}

func New() *Server {
	app := fiber.New()
	routes.Setup(app)
	return &Server{app: app}
}

func (server *Server) Start() {
	server.app.Listen(":80")
}

func (server *Server) Stop() {
	server.app.Shutdown()
}
