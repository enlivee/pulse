package repository

import (
	"sync"
	"github.com/enlivee/pulse/internal/model"
	"errors"
)

type InMemoryRepository struct {
	monitors map[string]*model.Monitor
	checks   map[string][]*model.Check
	mutex    sync.RWMutex
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		monitors: make(map[string]*model.Monitor),
		checks:   make(map[string][]*model.Check),
	}
}

func (r *InMemoryRepository) CreateMonitor(monitor *model.Monitor) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.monitors[monitor.ID] = monitor
	return nil
}

func (r *InMemoryRepository) GetMonitor(id string) (*model.Monitor, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	monitor, ok := r.monitors[id]
	if !ok {
		return nil, errors.New("Monitor not found")
	}
	copy := *monitor
	return &copy, nil
}

func (r *InMemoryRepository) GetMonitors() ([]*model.Monitor, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	monitors := make([]*model.Monitor, 0, len(r.monitors))
	for _, monitor := range r.monitors {
		copy := *monitor
		monitors = append(monitors, &copy)
	}
	return monitors, nil
}

func (r *InMemoryRepository) DeleteMonitor(id string) error{
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.monitors, id)
	delete(r.checks, id)
	return nil
}

func (r *InMemoryRepository) CreateCheck(check *model.Check) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if _, ok := r.monitors[check.MonitorID]; !ok {
		return errors.New("Monitor not found")
	}
	r.checks[check.MonitorID] = append(r.checks[check.MonitorID], check)
	return nil
}

func (r *InMemoryRepository) GetHistory(monitorID string) ([]*model.Check, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	history := make([]*model.Check, 0, len(r.checks[monitorID]))
	if _, ok := r.monitors[monitorID]; !ok {
		return nil, errors.New("Monitor not found")
	}
	for _, check := range r.checks[monitorID] {
		copy := *check
		history = append(history, &copy)
	}
	return history, nil
}