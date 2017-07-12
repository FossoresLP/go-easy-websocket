package websocket

import (
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/fossoreslp/go-jwt-ed25519"
	uuid "github.com/fossoreslp/go.uuid"
	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{}

// UpgradeHandler upgrades http requests to ws and starts a goroutine for handling ws messages
func (h *Handler) UpgradeHandler(w http.ResponseWriter, r *http.Request) {
	headers := strings.Split(r.Header.Get("Sec-Websocket-Protocol"), ",")
	auth, err := jwt.FromString(headers[0])
	if err != nil {
		w.WriteHeader(403)
		fmt.Println("JWT decoding")
		return
	}
	if !auth.Valid {
		w.WriteHeader(403)
		fmt.Println("JWT invalid")
		return
	}
	subject, err := uuid.FromString(auth.Content.Sub)
	if err != nil {
		w.WriteHeader(403)
		fmt.Println("JWT sub invalid")
		return
	}
	fmt.Println(subject.String())

	sessionid := uuid.NewV4()
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
