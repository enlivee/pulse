package repository

import "github.com/enlivee/pulse/internal/model"

type Repository interface {
	CreateMonitor(monitor *model.Monitor) error
	GetMonitor(id string) (*model.Monitor, error)
	GetMonitors() ([]*model.Monitor, error)
	DeleteMonitor(id string) error

	CreateCheck(check *model.Check) error
	GetHistory(monitorID string) ([]*model.Check, error)
}