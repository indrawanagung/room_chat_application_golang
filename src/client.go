package src

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
	"time"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}

type Message struct {
	Content  string `json:"content"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 4 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 6 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 4 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func (c *Client) writeMessage() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Message:
			if !ok {
				return
			}
			err := c.Conn.WriteJSON(message.Content)
			if err != nil {
				log.Error(err.Error())
			}

		case <-ticker.C:
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Error(err)
			}
			err = c.Conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Error(err.Error())
			}

			err = c.Conn.SetWriteDeadline(time.Time{})
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		err := c.Conn.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Error(err.Error())
	}

	c.Conn.SetPongHandler(func(string) error {
		err = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			log.Error(err.Error())
		}
		return nil
	})
	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Info(err.Error())
			}
			return
		}

		msg := &Message{
			Content:  string(m),
			RoomID:   c.RoomID,
			Username: c.Username,
		}
		hub.Broadcast <- msg
	}
}
