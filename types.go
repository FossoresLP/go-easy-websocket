package websocket

import "github.com/fossoreslp/go-uuid-v4"

// HandleFunc is a type used to store handle functions for ws commands
// Handle functions take the message as a byte slice and the auth token as a string and may return a response as a byte slice as well as an error
type HandleFunc func([]byte, string) ([]byte, error)

// channel channel and listening clients
type channel struct {
	send      chan []byte
	listeners []uuid.UUID
}

// Handler for a single websocket endpoint
type Handler struct {
	handlers      map[string]HandleFunc
	writeChannels map[uuid.UUID]chan []byte
	channels      map[string]*channel
}

// NewHandler creates a new Handler
func NewHandler() *Handler {
	return &Handler{
		make(map[string]HandleFunc),
		make(map[uuid.UUID]chan []byte),
		make(map[string]*channel),
	}
}
