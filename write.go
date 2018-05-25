package websocket

import (
	"errors"

	"github.com/fossoreslp/go-uuid-v4"
	ws "github.com/gorilla/websocket"
)

func (h *Handler) writerRoutine(conn *ws.Conn, userid uuid.UUID) {
	defer conn.Close()
	channel, ok := h.writeChannels[userid]
	if !ok {
		return
	}
	for {
		err := conn.WriteMessage(ws.TextMessage, <-channel)
		if err != nil {
			break
		}
	}
}

func (h *Handler) channelRoutine(channel string) {
	if c, ok := h.listeners[channel]; ok {
		for {
			msg := <-c.channel
			for _, listener := range c.listeners {
				err := h.WriteToClient(listener, msg)
				if err != nil {
					h.unregisterAsListener(listener, channel)
				}
			}
		}
	}
}

// WriteToChannel allows to write to all clients listening to a specific channel
func (h *Handler) WriteToChannel(channel string, data []byte) error {
	if c, ok := h.listeners[channel]; ok {
		c.channel <- data
		return nil
	}
	return errors.New("websocket: no such channel")
}

// WriteToClient allows to write to a single specific client
func (h *Handler) WriteToClient(user uuid.UUID, data []byte) error {
	if c, ok := h.writeChannels[user]; ok {
		c <- data
		return nil
	}
	return errors.New("websocket: client not found")
}
