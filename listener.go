package websocket

import (
	"errors"

	"github.com/fossoreslp/go-uuid-v4"
)

// RegisterListenChannel registers a channel for writing to clients listening on it
func (h *Handler) RegisterListenChannel(name string) error {
	if _, ok := h.listeners[name]; ok {
		return errors.New("channel already exists")
	}
	h.listeners[name] = &Listener{
		make(chan []byte, 8),
		make([]uuid.UUID, 1),
	}
	go h.channelRoutine(name)
	return nil
}

func (h *Handler) registerAsListener(id uuid.UUID, name string) error {
	if l, ok := h.listeners[name]; ok {
		l.listeners = append(l.listeners, id)
		return nil
	}
	return errors.New("websocket: channel does not exist")
}

func (h *Handler) unregisterAsListener(rmid uuid.UUID, name string) {
	for i, id := range h.listeners[name].listeners {
		if id == rmid {
			h.listeners[name].listeners = append(h.listeners[name].listeners[:i], h.listeners[name].listeners[i+1:]...)
		}
	}
}

func (h *Handler) unregisterListener(rmid uuid.UUID) {
	for name := range h.listeners {
		h.unregisterAsListener(rmid, name)
	}
}
