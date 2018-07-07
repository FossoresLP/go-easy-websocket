package websocket

import (
	"errors"
	"strings"

	"github.com/fossoreslp/go-uuid-v4"
)

// websocket command
var cmdWebSocket = []byte("websocket")

// HandleFunc is a type used to store handle functions for ws commands.
// Handle functions take the message as a byte slice and the auth token as a string and may return a message that will be submitted to the client or nil if no response is necessary.
type HandleFunc func([]byte, string) *Message

// channel stores a channel used to buffer the messsages as well as a slice containing the session ids of all listeners.
type channel struct {
	send      chan *Message
	listeners []uuid.UUID
}

/*Handler is the base type of a websocket endpoint.

It stores all relevant connections and is used to manage command handlers and channels.

The only public field is ValidateFunction which stores a function that is used to validate the users auth token.

Disabling authentication is currently not supported but you can simply supply a validation function that returns nil in all cases.
 func(_ string) error {
 	return nil
 }
Please make sure to set the auth cookie anyway as it is required for the connection to be accepted.*/
type Handler struct {
	ValidateFunction func(string) error // ValidateFunction is a function that validates the auth token and returns an error if it is invalid
	handlers         map[string]HandleFunc
	writeChannels    map[uuid.UUID]chan []byte
	channels         map[string]*channel
}

// NewHandler creates a new Handler and returns a pointer to it.
func NewHandler() *Handler {
	return &Handler{
		handlers:      make(map[string]HandleFunc),
		writeChannels: make(map[uuid.UUID]chan []byte),
		channels:      make(map[string]*channel),
	}
}

// Message is the type used to handle websocket messages.
// It is meant to be initialized using
//	NewMessage(command string, data []byte)
// to ensure only valid messages are created in the first place.
// For that reason the fields are not exported.
type Message struct {
	command []byte
	content []byte
}

/*NewMessage creates a new message from a command string and content submitted as a byte slice.
In case no content should be sent, please use nil instead of an empty slice.

There are some rules enforced by NewMessage():

- Commands may not be empty or longer than 255 characters

- Commands may not contain a colon

- Commands may not be "websocket" as this command is reserved*/
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
