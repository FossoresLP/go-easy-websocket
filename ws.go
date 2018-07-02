package websocket

import (
	"fmt"
	"net/http"

	"github.com/fossoreslp/go-uuid-v4"
	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	Subprotocols: []string{"cmd.fossores.de"},
}

// UpgradeHandler upgrades http requests to ws and starts a goroutine for handling ws messages
func (h *Handler) UpgradeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")
	if err != nil {
		w.WriteHeader(403)
		fmt.Fprintln(w, "Authentication failed")
		return
	}
	if h.ValidateFunction(cookie.Value) != nil {
		w.WriteHeader(403)
		fmt.Fprintln(w, "Authentication failed")
		return
	}

	sessionid, err := uuid.New()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(426)
		w.Header().Add("Upgrade", "WebSocket")
		return
	}
	h.writeChannels[sessionid] = make(chan []byte, 8)
	go h.handlerRoutine(conn, sessionid, cookie.Value)
	go h.writerRoutine(conn, sessionid)
}
