package websocket

import (
	"errors"
	"net/http"

	ws "github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

// HandleFunc is a type used to store handle functions for ws commands
type HandleFunc func([]byte) ([]byte, error)

var upgrader = ws.Upgrader{}

// Handler for a single websocket endpoint
type Handler struct {
	handlers map[string]HandleFunc
}

// NewHandler creates a new Handler
func NewHandler() *Handler {
	h := new(Handler)
	h.handlers = make(map[string]HandleFunc)
	return h
}

// UpgradeHandler upgrades http requests to wss and starts a goroutine for handling ws messages
func (h *Handler) UpgradeHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	go h.handlerRoutine(conn)
}

// handlerRoutine handles processing the recived messages and forwarding them to the defined handler functions
func (h *Handler) handlerRoutine(conn *ws.Conn) {
	defer conn.Close()
	if fnc, ok := h.handlers["open"]; ok {
		resp, err := fnc(nil)
		if err != nil {
			err = conn.WriteMessage(ws.TextMessage, []byte(err.Error()))
			if err != nil {
				return
			}
		}
		if resp != nil {
			err = conn.WriteMessage(ws.TextMessage, resp)
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
			err = conn.WriteMessage(ws.TextMessage, []byte("websocket: binary data not supported"))
			if err != nil {
				break
			}
		}
		cmd, data, err := parseMsgByte(msg)
		if err != nil {
			err = conn.WriteMessage(ws.TextMessage, []byte(err.Error()))
			if err != nil {
				break
			}
		}
		if fnc, ok := h.handlers[cmd]; ok {
			resp, err := fnc(data)
			if err != nil {
				err = conn.WriteMessage(ws.TextMessage, []byte("websocket: "+err.Error()))
				if err != nil {
					break
				}
			}
			if resp != nil {
				err = conn.WriteMessage(ws.TextMessage, resp)
				if err != nil {
					break
				}
			}
		} else {
			err = conn.WriteMessage(ws.TextMessage, []byte("websocket: command not supported by server"))
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
