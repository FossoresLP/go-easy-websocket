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
		resp, err := h.handlers[cmd](data)
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
	}
}

// Handle registers a handle function for a command
func (h *Handler) Handle(cmd string, action HandleFunc) {
	h.handlers[cmd] = action
}

/*func parseMsg(msg string) (cmd string, data string, err error) {
	split := strings.SplitN(msg, ": ", 2)
	if split[0] == "" {
		err = errors.New("websocket: no command defined")
		return
	}
	cmd = split[0]
	data = split[1]
	return
}*/

// Parse a recived message and return string and data
func parseMsgByte(msg []byte) (cmd string, data []byte, err error) {
	var i uint8
	for ; i <= 255; i++ {
		if msg[i] == ':' && msg[i+1] == ' ' {
			cmd = string(msg[0 : i-1])
			data = msg[i+2:]
			break
		}
	}
	if cmd == "" {
		err = errors.New("websocket: no command found")
	}
	return
}
