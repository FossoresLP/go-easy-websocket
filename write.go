package websocket

import (
	"errors"

	"github.com/fossoreslp/go-uuid-v4"
	ws "github.com/gorilla/websocket"
)

func (h *Handler) writerRoutine(conn *ws.Conn, userid uuid.UUID) {
	defer conn.Close()
	if channel, ok := h.writeChannels[userid]; ok {
		for {
			err := conn.WriteMessage(ws.TextMessage, <-channel)
			if err != nil {
				break
			}
		}
	}
}

func (h *Handler) channelRoutine(channel string) {
	if c, ok := h.channels[channel]; ok {
		for {
			msg := <-c.send
			for _, listener := range c.listeners {
				err := h.writeToClient(listener, msg.command, msg.content)
				if err != nil {
					h.unregisterAsListener(listener, channel)
				}
			}
		}
	}
}

// WriteToChannel allows to write to all clients listening to a specific channel
func (h *Handler) WriteToChannel(channel string, msg *Message) error {
	if c, ok := h.channels[channel]; ok {
		c.send <- msg
		return nil
	}
	return errors.New("channel does not exist")
}

// WriteToClient allows to write to a single specific client
func (h *Handler) WriteToClient(user uuid.UUID, msg *Message) error {
	if msg.command == nil {
		return errors.New("command may not be empty")
	}
	if len(msg.command) > 255 {
		return errors.New("command may not be longer than 255 characters")
	}
	return h.writeToClient(user, msg.command, msg.content)
}

func (h *Handler) writeToClient(user uuid.UUID, cmd, data []byte) error {
	if c, ok := h.writeChannels[user]; ok {
		prep := append(cmd, ':', ' ')
		c <- append(prep, data...)
		return nil
	}
	return errors.New("client not found")
}
