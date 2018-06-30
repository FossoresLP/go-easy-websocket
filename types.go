package websocket

import (
	"errors"
	"strings"

	"github.com/fossoreslp/go-uuid-v4"
)

// websocket command
var cmdWebSocket = []byte("websocket")

// HandleFunc is a type used to store handle functions for ws commands
// Handle functions take the message as a byte slice and the auth token as a string and may return a response as a byte slice as well as an error
type HandleFunc func([]byte, string) (*Message, error)

// channel channel and listening clients
type channel struct {
	send      chan *Message
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

// Message is a type that contains a websocket message
type Message struct {
	command []byte
	content []byte
}

// NewMessage creates a new message from a command string and the content as a byte slice
func NewMessage(cmd string, data []byte) (*Message, error) {
	if len(cmd) > 255 {
		return &Message{}, errors.New("command may not be longer than 255 characters")
	}
	if cmd == "" {
		return &Message{}, errors.New("command may not be emtpy")
	}
	if strings.ContainsRune(cmd, ':') {
		return &Message{}, errors.New("command may not contain a colon")
	}
	if cmd == "websocket" {
		return &Message{}, errors.New("command websocket is reserved")
	}
	return &Message{[]byte(cmd), data}, nil
}
