package main
import (
	"log"
	"github.com/gorilla/websocket"
	"net/http"
)

const (
	socketBufferSize	= 1024
	messageBufferSize	= 256
)

type room struct {
	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clients
	forward chan []byte
	// join is a channel for clients wishing to join the room
	join chan *client
	// leave is a channel for clients wishing to leave the room
	leave chan *client
	// client holds all current clients in this room
	clients map[*client]bool
}

func newRoom() *room {
	return &room{
		forward:	make(chan []byte),
		join:			make(chan *client),
		leave:		make(chan *client),
		clients:	make(map[*client]bool),
	}
}

func (r *room) Run() {
	for {
		select {
		case client := <-r.join:
			// joining
			r.clients[client] = true
		case client := <-r.join:
			// leaving
			close(client.send)
			delete(r.clients, client)
		case msg := <-r.forward:
			// forward message to all clients
		  for client := range r.clients {
				select {
				case client.send <- msg:
					// send the message
				default:
					// failed to send
					close(client.send)
					delete(r.clients, client)
				}
			}
		}
	}
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize: socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client {
		socket:	socket,
		send:		make(chan []byte, messageBufferSize),
		room:		r,
	}
	r.join <- client
	defer func() {r.leave <- client}()

	go client.write()
	client.read()
}