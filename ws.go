package websocket

import (
	"errors"
	"net/http"

	"strings"

	ws "github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

// HandleFunc is a type used to store handle functions for ws commands
type HandleFunc func(string) (string, error)

var upgrader = ws.Upgrader{}
var handlers = make(map[string]HandleFunc)

// HTTPHandler upgrades http requests to wss and starts a goroutine for handling ws messages
func HTTPHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	go WSSHandler(conn)
}

//WSSHandler handles messages from clients
func WSSHandler(conn *ws.Conn) {
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
		cmd, data, err := parseMsg(string(msg))
		if err != nil {
			err = conn.WriteMessage(ws.TextMessage, []byte(err.Error()))
			if err != nil {
				break
			}
		}
		resp, err := handlers[cmd](data)
		if err != nil {
			err = conn.WriteMessage(ws.TextMessage, []byte("websocket: "+err.Error()))
			if err != nil {
				break
			}
		}
		if resp != "" {
			err = conn.WriteMessage(ws.TextMessage, []byte(resp))
			if err != nil {
				break
			}
		}
	}
}

// Handle registers a handle function for a command
func Handle(cmd string, action HandleFunc) {
	handlers[cmd] = action
}

func parseMsg(msg string) (cmd string, data string, err error) {
	split := strings.SplitN(msg, ": ", 2)
	if split[0] == "" {
		err = errors.New("websocket: no command defined")
		return
	}
	cmd = split[0]
	data = split[1]
	return
}
