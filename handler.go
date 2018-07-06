package websocket

import (
	"errors"
	"reflect"
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
		if reflect.DeepEqual(msg.command, []byte("listen")) {
			if h.registerAsListener(sessionid, string(msg.content)) != nil {
				if h.writeToClient(sessionid, cmdWebSocket, []byte(err.Error())) != nil {
					break
				}
			}
		} else if fnc, ok := h.handlers[string(msg.command)]; ok {
			msg = fnc(msg.content, token)
			if msg.command != nil && msg.content != nil {
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
	var i uint8
	out := &Message{nil, nil}
	for ; i <= uint8(len(msg)-1) && i <= 255; i++ {
		if msg[i] == ':' && msg[i+1] == ' ' && len(msg[0:i]) > 0 {
			out.command = msg[0:i]
			if int(i+2) < len(msg) {
				out.content = msg[i+2:]
			}
			break
		}
	}
	return out
}
