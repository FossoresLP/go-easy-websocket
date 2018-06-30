package websocket

import (
	"errors"
	"strings"

	"github.com/fossoreslp/go-uuid-v4"
	ws "github.com/gorilla/websocket"
)

// handlerRoutine handles processing the recived messages and forwarding them to the defined handler functions
func (h *Handler) handlerRoutine(conn *ws.Conn, sessionid uuid.UUID, token string) {
	defer conn.Close()
	defer h.unregisterListener(sessionid)
	if fnc, ok := h.handlers["open"]; ok {
		msg, err := fnc([]byte(sessionid.String()), token)
		if err != nil {
			err = h.writeToClient(sessionid, cmdWebSocket, []byte(err.Error()))
			if err != nil {
				return
			}
		}
		if msg.command != nil && msg.content != nil {
			err = h.writeToClient(sessionid, msg.command, msg.content)
			if err != nil {
				return
			}
		}
	}
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		if msgType == ws.BinaryMessage {
			err = h.writeToClient(sessionid, cmdWebSocket, []byte("binary data not supported"))
			if err != nil {
				break
			}
		}
		cmd, data, err := parseMsgByte(msg)
		if err != nil {
			err = h.writeToClient(sessionid, cmdWebSocket, []byte(err.Error()))
			if err != nil {
				break
			}
		}
		if cmd == "listen" {
			err = h.registerAsListener(sessionid, string(data))
			if err != nil {
				err = h.writeToClient(sessionid, cmdWebSocket, []byte(err.Error()))
				if err != nil {
					break
				}
			}
		} else if fnc, ok := h.handlers[cmd]; ok {
			msg, err := fnc(data, token)
			if err != nil {
				err = h.writeToClient(sessionid, cmdWebSocket, []byte(err.Error()))
				if err != nil {
					break
				}
			}
			if msg.command != nil && msg.content != nil {
				err = h.writeToClient(sessionid, msg.command, msg.content)
				if err != nil {
					break
				}
			}
		} else {
			err = h.writeToClient(sessionid, cmdWebSocket, []byte("command not supported by server"))
			if err != nil {
				break
			}
		}
	}
}

// Handle registers a handle function for a command
func (h *Handler) Handle(cmd string, action HandleFunc) error {
	if len(cmd) < 255 {
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

// Parse a recived message and return string and data
func parseMsgByte(msg []byte) (cmd string, data []byte, err error) {
	var i uint8
	for ; i <= 255; i++ {
		if msg[i] == ':' && msg[i+1] == ' ' {
			cmd = string(msg[0:i])
			data = msg[i+2:]
			break
		}
	}
	if cmd == "" {
		err = errors.New("command not set")
	}
	return
}
