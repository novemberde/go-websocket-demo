package api

import "log"

// Hub Structure for Managing clients
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	send       chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run Start websocket server.
func (h *Hub) Run() {
	/**
	{
		event: "channel_message",
		channelId: body.channelId,
		name,
		content,
	}
	*/

	for {
		select {
		case client := <-h.register:
			// log.Println(client.)

			str := client.conn.RemoteAddr().String()
			client.id = getRandomString(16)
			log.Println(str)
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}

		}
	}
}
