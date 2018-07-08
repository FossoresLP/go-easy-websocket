package websocket

import (
	"errors"

	"github.com/fossoreslp/go-uuid-v4"
)

// RegisterListenChannel registers a channel for writing to clients listening on it
func (h *Handler) RegisterListenChannel(name string, validationFunc func(string) error) error {
	if _, ok := h.channels[name]; ok {
		return errors.New("channel already exists")
	}
	h.channels[name] = &channel{
		make(chan *Message, 8),
		make([]uuid.UUID, 0),
		validationFunc,
	}
	go h.channelRoutine(name)
	return nil
}

func (h *Handler) registerAsListener(id uuid.UUID, name string) error {
	if c, ok := h.channels[name]; ok {
		for _, lid := range c.listeners {
			if id == lid {
				return errors.New("already listening")
			}
		}
		c.listeners = append(c.listeners, id)
		return nil
	}
	return errors.New("channel does not exist")
}

func (h *Handler) unregisterAsListener(rmid uuid.UUID, name string) error {
	if c, ok := h.channels[name]; ok {
		for i, id := range c.listeners {
			if id == rmid {
				if len(c.listeners) <= 1 {
					c.listeners = []uuid.UUID{}
					break
				}
				if i == 0 {
					c.listeners = c.listeners[1:]
					break
				}
				if i == len(c.listeners)-1 {
					c.listeners = c.listeners[:i]
					break
				}
				c.listeners = append(c.listeners[:i], c.listeners[i+1:]...)
				break
			}
		}
		return nil
	}
	return errors.New("channel does not exist")
}

func (h *Handler) unregisterListener(rmid uuid.UUID) {
	for name := range h.channels {
		h.unregisterAsListener(rmid, name) // nolint: errcheck
	}
}
