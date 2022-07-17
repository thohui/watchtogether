package main

import (
	"github.com/thohui/watchtogether/server"
)

func main() {
	server := server.New()
	server.Start()
}
