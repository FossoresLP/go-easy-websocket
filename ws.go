package websocket

import (
	"fmt"
	"net/http"

	jwt "github.com/fossoreslp/go-jwt-ed25519"
	uuid "github.com/fossoreslp/go.uuid"
	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{}

// UpgradeHandler upgrades http requests to ws and starts a goroutine for handling ws messages
func (h *Handler) UpgradeHandler(w http.ResponseWriter, r *http.Request) {
	user, token, ok := r.BasicAuth()
	fmt.Println(ok)
	userid, err := uuid.FromString(user)
	if err != nil {
		w.WriteHeader(403)
		fmt.Println("UUID conversion")
		return
	}
	auth, err := jwt.Decode(token)
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
	if subject != userid {
		w.WriteHeader(403)
		fmt.Println("UUID != JWT sub")
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
