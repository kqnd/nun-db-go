package response

import (
	"strings"
)

type Handler struct {
	command  string
	args     []string
	watchers *map[string][]func(interface{})
	pendings *[]chan interface{}
}

func CreateHandler(watchers *map[string][]func(interface{}), pendings *[]chan interface{}) Handler {
	return Handler{
		watchers: watchers,
		pendings: pendings,
	}
}

func (h *Handler) SetPayload(response string) {
	args := splitResponse(response)
	if len(args) == 0 {
		h.command = ""
		h.args = nil
		return
	}
	h.command = args[0]
	if len(args) > 1 {
		h.args = args[1:]
	} else {
		h.args = nil
	}
}

func (h *Handler) GettingValues() {
	if h.command == "value" {
		value := h.args[0]

		pendings := *h.pendings
		if len(pendings) == 0 {
			return
		}

		ch := (pendings)[0]
		ch <- value
		close(ch)

		*h.pendings = pendings[1:]
	}
}

func (h *Handler) WatchingValues(strictQueueWatch bool) {
	if h.command == "changed" {
		key := h.args[0]
		value := h.args[1]
		if watchers, ok := (*h.watchers)[key]; ok {
			for _, cb := range watchers {
				if strictQueueWatch {
					cb(value)
				} else {
					go cb(value)
				}
			}
		}
	}
}

func hidePrivateKeys(keys []string) []string {
	newKeys := []string{}
	for _, key := range keys {
		if key != "$$token" && key != "$connections" {
			newKeys = append(newKeys, key)
		}
	}
	return newKeys
}

func (h *Handler) Keys() {
	if h.command == "keys" {
		var keys []string = strings.Split(h.args[0], ",")

		pendings := *h.pendings
		if len(pendings) == 0 {
			return
		}

		ch := (pendings)[0]
		ch <- hidePrivateKeys(keys)[1:]
		close(ch)

		*h.pendings = pendings[1:]
	}
}
