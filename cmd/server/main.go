package main

import (
	"url-shortener/configs"
	"url-shortener/internal/app"
)

func main() {
	config := configs.ReadConfig()
	server := app.New(config)

	server.Run()
}
