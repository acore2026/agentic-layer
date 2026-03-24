package events

import (
	"fmt"
	"log"
	"net/http"
)

// Broker handles multiple SSE client connections and broadcasts messages to them.
type Broker struct {
	// Notifier channel for broadcasting messages
	Notifier chan []byte

	// New clients wishing to connect
	newClients chan chan []byte

	// Clients wishing to disconnect
	closingClients chan chan []byte

	// Currently active clients
	clients map[chan []byte]bool
}

// NewBroker initializes and starts an SSE broker.
func NewBroker() *Broker {
	broker := &Broker{
		Notifier:       make(chan []byte, 1),
		newClients:      make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}

	// Start listening for client events and notifications
	go broker.listen()

	return broker
}

func (b *Broker) listen() {
	for {
		select {
		case s := <-b.newClients:
			// A new client has connected
			b.clients[s] = true
			log.Printf("SSE: Client connected. Total clients: %d", len(b.clients))

		case s := <-b.closingClients:
			// A client has disconnected
			delete(b.clients, s)
			log.Printf("SSE: Client disconnected. Total clients: %d", len(b.clients))

		case event := <-b.Notifier:
			// Broadcast the event to all connected clients
			for clientChan := range b.clients {
				select {
				case clientChan <- event:
					// Message sent
				default:
					// Slow client, skip or drop
					log.Println("SSE: Skipping slow client")
				}
			}
		}
	}
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Ensure the connection supports flushing
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Set mandatory SSE and CORS headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a message channel for this client
	messageChan := make(chan []byte)

	// Register this client
	b.newClients <- messageChan

	// Unregister on connection close
	defer func() {
		b.closingClients <- messageChan
	}()

	// Keep connection alive and wait for messages
	for {
		select {
		case msg := <-messageChan:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
