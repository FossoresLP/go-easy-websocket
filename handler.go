package websocket

import (
	"errors"

	"github.com/fossoreslp/go-uuid-v4"
	ws "github.com/gorilla/websocket"
)

// handlerRoutine handles processing the recived messages and forwarding them to the defined handler functions
func (h *Handler) handlerRoutine(conn *ws.Conn, userid uuid.UUID) {
	defer conn.Close()
	defer h.unregisterListener(userid)
	if fnc, ok := h.handlers["open"]; ok {
		resp, err := fnc(nil, userid)
		if err != nil {
			err = h.WriteToClient(userid, []byte(err.Error()))
			if err != nil {
				return
			}
		}
		if resp != nil {
			err = h.WriteToClient(userid, resp)
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
			err = h.WriteToClient(userid, []byte("websocket: binary data not supported"))
			if err != nil {
				break
			}
		}
		cmd, data, err := parseMsgByte(msg)
		if err != nil {
			err = h.WriteToClient(userid, []byte(err.Error()))
			if err != nil {
				break
			}
		}
		if cmd == "listen" {
			err = h.registerAsListener(userid, string(data))
			if err != nil {
				err = h.WriteToClient(userid, []byte(err.Error()))
				if err != nil {
					break
				}
			}
		} else if fnc, ok := h.handlers[cmd]; ok {
			resp, err := fnc(data, userid)
			if err != nil {
				err = h.WriteToClient(userid, []byte("websocket: "+err.Error()))
				if err != nil {
					break
				}
			}
			if resp != nil {
				err = h.WriteToClient(userid, resp)
				if err != nil {
					break
				}
			}
		} else {
			err = h.WriteToClient(userid, []byte("websocket: command not supported by server"))
			if err != nil {
				break
			}
		}
	}
}

// Handle registers a handle function for a command
func (h *Handler) Handle(cmd string, action HandleFunc) {
	h.handlers[cmd] = action
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
		err = errors.New("websocket: no command found")
	}
	return
}
