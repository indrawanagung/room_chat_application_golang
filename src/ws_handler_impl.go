package src

import (
	"fmt"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type WebsocketHandlerImpl struct {
	hub *Hub
}

func NewWebsocketHandler(h *Hub) WebsocketHandlerInterface {
	return &WebsocketHandlerImpl{
		hub: h,
	}
}

func (h *WebsocketHandlerImpl) CreateRoom(c *fiber.Ctx) error {
	req := new(CreateRoomReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).Send([]byte("error create room request"))
	}

	h.hub.Rooms[req.ID] = &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}

	return c.Status(200).JSON(req)
}

func (h *WebsocketHandlerImpl) JoinRoom(c *websocket.Conn) {
	roomID := c.Params("id")
	username := c.Query("username")

	cl := &Client{
		Conn:     c,
		Message:  make(chan *Message, 10),
		RoomID:   roomID,
		Username: username,
	}

	// _, ok := h.hub.Rooms[roomID]
	// if !ok {
	// 	c.WriteJSON("room id is not found")
	// 	c.Conn.Close()
	// 	return
	// }
	// _, ok = h.hub.Rooms[roomID].Clients[username]
	// if ok {
	// 	c.WriteJSON("username has already exist")
	// 	c.Conn.Close()
	// 	return
	// }

	m := &Message{
		Content:  fmt.Sprintf("%s has joined the room", username),
		RoomID:   roomID,
		Username: username,
	}

	h.hub.Register <- cl
	h.hub.Broadcast <- m
	go cl.writeMessage()
	cl.readMessage(h.hub)
}

func (h *WebsocketHandlerImpl) GetListRoom(c *fiber.Ctx) error {
	rooms := make([]RoomRes, 0)
	for _, r := range h.hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	return c.Status(200).JSON(rooms)
}

func (h *WebsocketHandlerImpl) GetListClient(c *fiber.Ctx) error {
	var clients []ClientRes
	roomId := c.Params("roomId")

	if _, ok := h.hub.Rooms[roomId]; !ok {
		clients = make([]ClientRes, 0)
		return c.Status(200).JSON(clients)
	}

	for _, c := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, ClientRes{
			Username: c.Username,
		})
	}

	return c.Status(200).JSON(clients)
}
