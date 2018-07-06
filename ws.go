/*Package websocket provides an easy, text-based way to build a websocket server.

It uses a command based message structure. The command can have a length between 1 and 255 characters and may contain all characters besides a colon although it is recommended to use short, lowercase, one word commands to keep the messages as short as possible.

The command is followed by a colon and a space after which the optional content is appended.

This format keeps all messages as human-readable as possible to allow for easy diagnosis of errors.

You can register a handler for a specific command which allows you to respond directly. This is the recommended way of responding to a command although it is also possible to store the users session id and use it to respond later. This might be useful when complex computations are necessary to generate the response and you don't want to block the read loop.

It is possible to address multiple clients at once through the use of channels. You have to register a channel name which clients can then use to register as listeners. A message submitted to a channel will be sent to all listening clients.*/
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
		fmt.Fprintln(w, "Authentication failed") // nolint: errcheck
		return
	}
	if h.ValidateFunction(cookie.Value) != nil {
		w.WriteHeader(403)
		fmt.Fprintln(w, "Authentication failed") // nolint: errcheck
		return
	}

	sessionid, err := uuid.New()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, "Server failed to initialize session") // nolint: errcheck
	}
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
