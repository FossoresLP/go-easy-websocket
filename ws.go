package websocket

import (
	"fmt"
	"net/http"

	"github.com/fossoreslp/go-uuid-v4"
	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{}

// ValidateFunction is a function that validates the auth token and returns an error if it is invalid
var ValidateFunction func(string) error

// UpgradeHandler upgrades http requests to ws and starts a goroutine for handling ws messages
func (h *Handler) UpgradeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")
	if err != nil {
		w.WriteHeader(403)
		fmt.Println("Authentication failed")
		return
	}
	if !cookie.HttpOnly || !cookie.Secure {
		w.WriteHeader(403)
		fmt.Println("Authentication failed")
		return
	}
	if ValidateFunction(cookie.Value) != nil {
		w.WriteHeader(403)
		fmt.Println("Authentication failed")
		return
	}

	sessionid, err := uuid.New()
	upgrader.Subprotocols = []string{"cmd.fossores.de"}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(426)
		w.Header().Add("Upgrade", "WebSocket")
		return
	}
	h.writeChannels[sessionid] = make(chan []byte, 8)
	go h.handlerRoutine(conn, sessionid)
	go h.writerRoutine(conn, sessionid)
}
