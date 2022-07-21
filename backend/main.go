package main

import (
	_ "github.com/joho/godotenv/autoload"

	"github.com/thohui/watchtogether/server"
)

func main() {
	server := server.New()
	server.Start()
}
