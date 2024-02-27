package main

import (
	"gochat/src"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	app := fiber.New()

	hub := src.NewHub()
	wsHandler := src.NewWebsocketHandler(hub)
	go hub.Run()

	app.Post("/ws/createRoom", wsHandler.CreateRoom)
	app.Get("/ws/joinRoom/:id", websocket.New(wsHandler.JoinRoom))
	app.Get("/ws/getRooms", wsHandler.GetListRoom)
	app.Get("/ws/getClients/:roomId", wsHandler.GetListClient)

	log.Fatal(app.Listen(":3000"))
}
