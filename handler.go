package websocket

import (
	"bytes"
	"errors"
	"strings"

	"github.com/fossoreslp/go-uuid-v4"
	ws "github.com/gorilla/websocket"
)

// handlerRoutine handles processing the recived messages and forwarding them to the defined handler functions
func (h *Handler) handlerRoutine(conn *ws.Conn, sessionid uuid.UUID, token string) {
	defer conn.Close() // nolint: errcheck
	defer h.unregisterListener(sessionid)
	if fnc, ok := h.handlers["open"]; ok {
		msg := fnc([]byte(sessionid.String()), token)
		if msg.command != nil && msg.content != nil {
			if h.writeToClient(sessionid, msg.command, msg.content) != nil {
				return
			}
		}
	}
	for {
		_, rawMsg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		msg := parseMessage(rawMsg)
		if bytes.Equal(msg.command, []byte("listen")) {
			if err := h.registerAsListener(sessionid, string(msg.content)); err != nil {
				if h.writeToClient(sessionid, cmdWebSocket, []byte(err.Error())) != nil {
					break
				}
			}
		} else if fnc, ok := h.handlers[string(msg.command)]; ok {
			msg = fnc(msg.content, token)
			if msg != nil && msg.command != nil && msg.content != nil {
				if h.writeToClient(sessionid, msg.command, msg.content) != nil {
					break
				}
			}
		} else {
			if h.writeToClient(sessionid, cmdWebSocket, []byte("command not supported by server")) != nil {
				break
			}
		}
	}
}

// Handle registers a handle function for a command
func (h *Handler) Handle(cmd string, action HandleFunc) error {
	if len(cmd) > 255 {
		return errors.New("command may not be longer than 255 characters")
	}
	if strings.ContainsRune(cmd, ':') {
		return errors.New("command may not contain a colon")
	}
	if cmd == "websocket" {
		return errors.New("command websocket is reserved")
	}
	if _, ok := h.handlers[cmd]; ok {
		return errors.New("command already exists")
	}
	h.handlers[cmd] = action
	return nil
}

// parseMessage returns a Message pointer
func parseMessage(msg []byte) *Message {
	out := &Message{nil, nil}
	var i int
	if len(msg) < 256 {
		i = bytes.Index(msg, []byte(": "))
	} else {
		i = bytes.Index(msg[:256], []byte(": "))
	}
	if i >= 1 {
		out.command = msg[:i]
		if i+2 < len(msg) {
			out.content = msg[i+2:]
		}
	}
	return out
}
