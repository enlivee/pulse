package repository_test

import (
	"testing"

	"github.com/enlivee/pulse/internal/model"
	"github.com/enlivee/pulse/internal/repository"
)

func TestInMemoryRepository_CreateMonitor(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	monitor := &model.Monitor{ID: "1", URL: "Test Monitor"}

	err := repo.CreateMonitor(monitor)
	if err != nil {
		t.Fatalf("CreateMonitor failed: %v", err)
	}
}

func TestInMemoryRepository_GetMonitor(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	monitor := &model.Monitor{ID: "1", URL: "Test Monitor"}

	err := repo.CreateMonitor(monitor)
	if err != nil {
		t.Fatalf("CreateMonitor failed: %v", err)
	}

	retrievedMonitor, err := repo.GetMonitor("1")
	if err != nil {
		t.Fatalf("GetMonitor failed: %v", err)
	}
	if retrievedMonitor.ID != monitor.ID || retrievedMonitor.URL != monitor.URL {
		t.Fatalf("GetMonitor returned incorrect data: got %v, want %v", retrievedMonitor, monitor)
	}

	retrievedMonitor.URL = "modified"
	newRetrievedMonitor, err := repo.GetMonitor("1")
	if err != nil {
		t.Fatalf("GetMonitor failed: %v", err)
	}
	if retrievedMonitor.URL == newRetrievedMonitor.URL {
		t.Fatalf("GetMonitor returned a reference to the original monitor, expected a copy")
	}
}

func TestInMemoryRepository_GetMonitors(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	monitor := &model.Monitor{ID: "1", URL: "Test Monitor"}

	err := repo.CreateMonitor(monitor)
	if err != nil {
		t.Fatalf("CreateMonitor failed: %v", err)
	}

	monitors, err := repo.GetMonitors()
	if err != nil {
		t.Fatalf("GetMonitors failed: %v", err)
	}
	if len(monitors) != 1 || monitors[0].ID != monitor.ID {
		t.Fatalf("GetMonitors returned incorrect data: got %v, want %v", monitors, []*model.Monitor{monitor})
	}
}

func TestInMemoryRepository_CreateCheck(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	monitor := &model.Monitor{ID: "1", URL: "Test Monitor"}

	err := repo.CreateMonitor(monitor)
	if err != nil {
		t.Fatalf("CreateMonitor failed: %v", err)
	}

	check := &model.Check{MonitorID: "1", StatusCode: 200}
	err = repo.CreateCheck(check)
	if err != nil {
		t.Fatalf("CreateCheck failed: %v", err)
	}

	notExistingCheck := &model.Check{MonitorID: "non-existing", StatusCode: 200}
	err = repo.CreateCheck(notExistingCheck)
	if err == nil {
		t.Fatalf("CreateCheck for non-existing monitor failed: %v", err)
	}
}

func TestInMemoryRepository_GetHistory(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	monitor := &model.Monitor{ID: "1", URL: "Test Monitor"}

	err := repo.CreateMonitor(monitor)
	if err != nil {
		t.Fatalf("CreateMonitor failed: %v", err)
	}

	check := &model.Check{MonitorID: "1", StatusCode: 200}
	err = repo.CreateCheck(check)
	if err != nil {
		t.Fatalf("CreateCheck failed: %v", err)
	}

	history, err := repo.GetHistory("1")
	if err != nil {
		t.Fatalf("GetHistory failed: %v", err)
	}
	if len(history) != 1 || history[0].StatusCode != 200 {
		t.Fatalf("GetHistory returned incorrect data: got %v, want %v", history, []*model.Check{check})
	}

	history[0].StatusCode = 404
	newHistory, err := repo.GetHistory("1")
	if err != nil {
		t.Fatalf("GetHistory failed: %v", err)
	}

	if newHistory[0].StatusCode != 200 {
		t.Fatalf("GetHistory returned a reference to the original check")
	}
}

func TestInMemoryRepository_DeleteMonitor(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	monitor := &model.Monitor{ID: "1", URL: "Test Monitor"}

	err := repo.CreateMonitor(monitor)
	if err != nil {
		t.Fatalf("CreateMonitor failed: %v", err)
	}

	check := &model.Check{MonitorID: "1", StatusCode: 200}
	err = repo.CreateCheck(check)
	if err != nil {
		t.Fatalf("CreateCheck failed: %v", err)
	}

	err = repo.DeleteMonitor("1")
	if err != nil {
		t.Fatalf("DeleteMonitor failed: %v", err)
	}

	_, err = repo.GetMonitor("1")
	if err == nil {
		t.Fatalf("GetMonitor should have failed for deleted monitor")
	}

	_, err = repo.GetHistory("1")
	if err == nil {
		t.Fatalf("GetHistory should have failed for deleted monitor")
	}
}
