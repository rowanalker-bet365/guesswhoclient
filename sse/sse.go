package sse

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

// Broker maintains the set of active clients and broadcasts messages to the clients.
type Broker struct {
	clients    map[chan string]bool
	newClients chan chan string
	defunctClients chan chan string
	messages   chan string
	mu         sync.Mutex
}

// NewBroker creates a new Broker.
func NewBroker() *Broker {
	b := &Broker{
		clients:    make(map[chan string]bool),
		newClients: make(chan (chan string)),
		defunctClients: make(chan (chan string)),
		messages:   make(chan string),
	}
	go b.listen()
	return b
}

// listen runs in a goroutine and manages the clients and messages.
func (b *Broker) listen() {
	for {
		select {
		case s := <-b.newClients:
			b.mu.Lock()
			b.clients[s] = true
			b.mu.Unlock()
			log.Println("Client added. Total clients:", len(b.clients))
		case s := <-b.defunctClients:
			b.mu.Lock()
			delete(b.clients, s)
			b.mu.Unlock()
			close(s)
			log.Println("Client removed. Total clients:", len(b.clients))
		case msg := <-b.messages:
			b.mu.Lock()
			for s := range b.clients {
				s <- msg
			}
			b.mu.Unlock()
			log.Println("Broadcasted message to all clients:", msg)
		}
	}
}

// ServeHTTP handles HTTP requests for the SSE endpoint.
func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	messageChan := make(chan string)
	b.newClients <- messageChan

	defer func() {
		b.defunctClients <- messageChan
	}()

	notify := r.Context().Done()
	go func() {
		<-notify
		b.defunctClients <- messageChan
	}()

	for {
		fmt.Fprintf(w, "data: %s\n\n", <-messageChan)
		flusher.Flush()
	}
}

// Broadcast sends a message to all clients.
func (b *Broker) Broadcast(msg string) {
	b.messages <- msg
}