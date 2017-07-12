package websocket

import (
	"net/http"

	uuid "github.com/fossoreslp/go.uuid"
	ws "github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

var upgrader = ws.Upgrader{}

// UpgradeHandler upgrades http requests to ws and starts a goroutine for handling ws messages
func (h *Handler) UpgradeHandler(w http.ResponseWriter, r *http.Request, params ...httprouter.Params) {
	if err != nil {
		w.WriteHeader(500)
		return
	}
	sessionid := uuid.NewV4()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(426)
		w.Header().Add("Upgrade", "WebSocket")
		return
	}
	h.writeChannels[userid] = make(chan []byte, 8)
	go h.handlerRoutine(conn, sessionid)
	go h.writerRoutine(conn, sessionid)
}
