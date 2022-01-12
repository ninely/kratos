package health

import (
	"sync"
)

type watch func(status Status)
type Health struct {
	statusMap map[string]Status
	mutex     sync.RWMutex
	opts      options
	watchers  []watch
}
type Option func(*options)

type options struct{}

func New(opts ...Option) *Health {
	h := &Health{
		statusMap: make(map[string]Status),
		mutex:     sync.RWMutex{},
	}
	option := options{}
	for _, o := range opts {
		o(&option)
	}
	return h
}

func (h *Health) SetStatus(service string, status Status) {
	h.mutex.Lock()
	h.statusMap[service] = status
	h.mutex.Unlock()
	go func() {
		for _, w := range h.watchers {
			w(status)
		}
	}()
}

func (h *Health) GetStatus(service string) (status Status, ok bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	status, ok = h.statusMap[service]
	return status, ok
}

func (h *Health) Watch(service string, f watch) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	h.watchers = append(h.watchers, f)
}
