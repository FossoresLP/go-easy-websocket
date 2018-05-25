package websocket

import "github.com/fossoreslp/go-uuid-v4"

// HandleFunc is a type used to store handle functions for ws commands
type HandleFunc func([]byte, uuid.UUID) ([]byte, error)

// Listener channel and listening clients
type Listener struct {
	channel   chan []byte
	listeners []uuid.UUID
}

// Handler for a single websocket endpoint
type Handler struct {
	handlers      map[string]HandleFunc
	writeChannels map[uuid.UUID]chan []byte
	listeners     map[string]*Listener
}

// NewHandler creates a new Handler
func NewHandler() *Handler {
	h := new(Handler)
	h.handlers = make(map[string]HandleFunc)
	h.writeChannels = make(map[uuid.UUID]chan []byte)
	h.listeners = make(map[string]*Listener)
	return h
}
