package src

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type WebsocketHandlerInterface interface {
	CreateRoom(c *fiber.Ctx) error
	JoinRoom(c *websocket.Conn)
	GetListRoom(c *fiber.Ctx) error
	GetListClient(c *fiber.Ctx) error
}
