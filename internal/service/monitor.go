package service

import (
	"errors"
	"github.com/enlivee/pulse/internal/model"
	"github.com/enlivee/pulse/internal/repository"
	"github.com/google/uuid"
	"time"
)

type MonitorService struct {
	rep repository.Repository
}

func NewMonitorService(rep repository.Repository) *MonitorService {
	return &MonitorService{rep: rep}
}

func generateID() string {
	return uuid.New().String()
}

func (s *MonitorService) CreateMonitor(URL string, interval time.Duration) (*model.Monitor, error) {
	if URL == "" {
		return nil, errors.New("URL cannot be empty")
	}
	monitor := &model.Monitor{
		ID:       generateID(),
		URL:      URL,
		Interval: interval,
	}

	err := s.rep.CreateMonitor(monitor)
	if err != nil {
		return nil, err
	}
	return monitor, nil
}

func (s *MonitorService) GetMonitor(id string) (*model.Monitor, error) {
	return s.rep.GetMonitor(id)
}

func (s *MonitorService) GetMonitors() ([]*model.Monitor, error) {
	return s.rep.GetMonitors()
}

func (s *MonitorService) DeleteMonitor(id string) error {
	return s.rep.DeleteMonitor(id)
}

func (s *MonitorService) CreateCheck(monitorID string, statusCode int, responseTime time.Duration) (*model.Check, error) {
	check := &model.Check{
		MonitorID:  monitorID,
		Time:       time.Now(),
		StatusCode: statusCode,
		Latency:    responseTime,
		Success:    statusCode >= 200 && statusCode < 300,
	}
	return check, s.rep.CreateCheck(check)
}

func (s *MonitorService) GetHistory(monitorID string) ([]*model.Check, error) {
	return s.rep.GetHistory(monitorID)
}
