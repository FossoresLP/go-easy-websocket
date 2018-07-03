package websocket

import (
	"errors"

	"github.com/fossoreslp/go-uuid-v4"
	ws "github.com/gorilla/websocket"
)

// writerRoutine is the goroutine spawned to send all messages that are queued for a specific client.
// It will check if a channel exists for the messages to send and then indefinitely loop over the incoming messages on that channel and send those to the client.
// The loop will exit when a write fails. This should only ever happen if the client disconnected.
// This goroutine will close the connection to the client upon exiting.
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

// channelRoutine is the goroutine spawned to handle all messages that are queued for a specific channel.
// It will check if the channel exists and then indefinitely loop over the incoming messages trying to send them to all registered listeners.
// If writing to a listener fails which will only ever happen when that listener is no longer connected, the session id removed as a listener.
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

// WriteToChannel sends a message to all clients listening to a specific channel.
// It takes the channel name and a pointer to a message as arguments.
// It will fail if the command is nil or longer than 255 characters or if the channel does not exist.
func (h *Handler) WriteToChannel(channel string, msg *Message) error {
	if msg.command == nil {
		return errors.New("command may not be empty")
	}
	if len(msg.command) > 255 {
		return errors.New("command may not be longer than 255 characters")
	}
	if c, ok := h.channels[channel]; ok {
		c.send <- msg
		return nil
	}
	return errors.New("channel does not exist")
}

// WriteToClient sends a message to a specific client.
// It takes the users session id as an UUID and a pointer to a message as arguments.
// It will fail if the command is nil or longer than 255 characters or if the session does not exist.
func (h *Handler) WriteToClient(user uuid.UUID, msg *Message) error {
	if msg.command == nil {
		return errors.New("command may not be empty")
	}
	if len(msg.command) > 255 {
		return errors.New("command may not be longer than 255 characters")
	}
	return h.writeToClient(user, msg.command, msg.content)
}

// writeToClient is the underlying function that is used send messages the individual clients.
//It takes the userid, command and message.
// These are then combined into the correct message format and passed to the send channel.
func (h *Handler) writeToClient(user uuid.UUID, cmd, data []byte) error {
	if c, ok := h.writeChannels[user]; ok {
		prep := append(cmd, ':', ' ')
		c <- append(prep, data...)
		return nil
	}
	return errors.New("client not found")
}
