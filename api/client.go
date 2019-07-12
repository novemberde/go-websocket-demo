package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketEvent struct {
	// Required
	Event string `json:"event"`

	//
	Message string `json:"message"`
	// FIXME: SenderID should be set on server using token. UserID.
	SenderID int    `json:"sender_id"`
	RoomID   string `json:"room_id"`
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10 // Must be less than pongWait
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// FIXME: Check origins on production
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
	hub   *Hub
	conn  *websocket.Conn
	send  chan []byte // Buffered channel of outbound messages.
	id    string
	rooms []string
}

//
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		// messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Fatalln(err)
			if err != nil {
				log.Fatalln(err)
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				break
			}
		}
		// Trim
		// message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		// message = bytes.TrimSpace(message)

		// Broadcate a messsage
		// c.hub.broadcast <- message

		// Create room
		var e *WebSocketEvent
		err = json.Unmarshal(message, &e)
		if err != nil {
			log.Fatalln(err)
			break
		}

		switch e.Event {
		case "connect":
			// {"event": "connect"}
			log.Println("connect")
		case "disconnect":
			log.Println("disconnect")
		case "room":
			// {"event": "room"}
			log.Println(c.rooms)
		case "room/join":
			// {"event": "room/join", "room_id": "BpLnfgDsc3WD9F3q"}
			// rooms
			// joinRoom(c, room)
			for _, r := range rooms {
				if r.RoomID == e.RoomID {
					r.Users = append(r.Users, c)
					log.Println("RoomID", r.RoomID)
					log.Println("users: ", len(r.Users))
				}
			}
		case "room/create":
			// {"event": "room/create"}
			log.Println("room/create")
			room := createRoom()
			c.rooms = append(c.rooms, room.RoomID)
			joinRoom(c, room)
			log.Println("ID: ", c.id)
		case "room/msg/send":
			log.Println("room/msg/send")

		case "dm/send":
			log.Println("dm/send")

		default:
		}

		// Send users a event
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			// TODO: 채팅방에 보내야함.
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Fatalln(err)
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(h *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// client := &Client{hub: h, conn: conn, send: make(chan []byte, 256)}
	client := &Client{hub: h, conn: conn, send: make(chan []byte)}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
