package websocket

import (
	"errors"

	"github.com/fossoreslp/go-uuid-v4"
)

// RegisterListenChannel registers a channel for writing to clients listening on it
func (h *Handler) RegisterListenChannel(name string) error {
	if _, ok := h.channels[name]; ok {
		return errors.New("channel already exists")
	}
	h.channels[name] = &channel{
		make(chan *Message, 8),
		make([]uuid.UUID, 1),
	}
	go h.channelRoutine(name)
	return nil
}

func (h *Handler) registerAsListener(id uuid.UUID, name string) error {
	if c, ok := h.channels[name]; ok {
		c.listeners = append(c.listeners, id)
		return nil
	}
	return errors.New("channel does not exist")
}

func (h *Handler) unregisterAsListener(rmid uuid.UUID, name string) {
	for i, id := range h.channels[name].listeners {
		if id == rmid {
			h.channels[name].listeners = append(h.channels[name].listeners[:i], h.channels[name].listeners[i+1:]...)
		}
	}
}

func (h *Handler) unregisterListener(rmid uuid.UUID) {
	for name := range h.channels {
		h.unregisterAsListener(rmid, name)
	}
}
