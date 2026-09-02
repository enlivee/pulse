package repository

import (
	"sync"
	"github.com/enlivee/pulse/internal/model"
)

type InMemoryRepository struct {
	monitors map[string]*model.Monitor
	checks   map[string][]*model.Check
	mutex    sync.Mutex
}