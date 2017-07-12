package websocket

import (
	"net/http"

	uuid "github.com/fossoreslp/go.uuid"
	ws "github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

var upgrader = ws.Upgrader{}

// UpgradeHandler upgrades http requests to wss and starts a goroutine for handling ws messages
// UpgradeHandler upgrades http requests to ws and starts a goroutine for handling ws messages
func (h *Handler) UpgradeHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userid, err := uuid.FromString(params.ByName("uuid"))
	if err != nil {
		w.WriteHeader(500)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(426)
		w.Header().Add("Upgrade", "WebSocket")
		return
	}
	h.writeChannels[userid] = make(chan []byte, 8)
	go h.handlerRoutine(conn, userid)
	go h.writerRoutine(conn, userid)
}
