package service_test

import (
	"strings"
	"testing"
	"time"

	"github.com/enlivee/pulse/internal/model"
	"github.com/enlivee/pulse/internal/repository"
	"github.com/enlivee/pulse/internal/service"
	"github.com/google/uuid"
)

func TestMonitorService_CreateMonitor(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	svc := service.NewMonitorService(repo)

	t.Run("ошибка при пустом URL", func(t *testing.T) {
		monitor, err := svc.CreateMonitor("", time.Minute)

		if err == nil {
			t.Fatal("ожидалась ошибка, получено nil")
		}
		if !strings.Contains(err.Error(), "URL cannot be empty") {
			t.Errorf("неверное сообщение ошибки: got %q", err.Error())
		}
		if monitor != nil {
			t.Error("ожидалось, что monitor будет nil")
		}
	})

	t.Run("успешное создание, генерация ID и сохранение", func(t *testing.T) {
		url := "https://example.com"
		interval := 30 * time.Second

		monitor, err := svc.CreateMonitor(url, interval)
		if err != nil {
			t.Fatalf("неожиданная ошибка: %v", err)
		}

		if monitor.URL != url {
			t.Errorf("неверный URL: got %q, want %q", monitor.URL, url)
		}
		if monitor.Interval != interval {
			t.Errorf("неверный Interval: got %v, want %v", monitor.Interval, interval)
		}

		// Проверка генерации валидного UUID
		if _, err := uuid.Parse(monitor.ID); err != nil {
			t.Errorf("ID не является валидным UUID: %v", err)
		}

		// Проверка, что данные действительно сохранены в repository
		savedMonitor, err := repo.GetMonitor(monitor.ID)
		if err != nil {
			t.Fatalf("монитор не найден в репозитории: %v", err)
		}
		if savedMonitor.ID != monitor.ID {
			t.Error("монитор не был корректно сохранен в репозитории")
		}
	})
}

func TestMonitorService_GetMonitor(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	svc := service.NewMonitorService(repo)

	// Создаем тестовый монитор напрямую через репозиторий для проверки чтения
	testMonitor := &model.Monitor{ID: "test-id-123", URL: "https://test.com", Interval: time.Minute}
	_ = repo.CreateMonitor(testMonitor)

	t.Run("успешное получение", func(t *testing.T) {
		m, err := svc.GetMonitor("test-id-123")
		if err != nil {
			t.Fatalf("неожиданная ошибка: %v", err)
		}
		if m.ID != testMonitor.ID || m.URL != testMonitor.URL {
			t.Errorf("неверные данные: got %+v, want %+v", m, testMonitor)
		}
	})

	t.Run("ошибка при несуществующем ID", func(t *testing.T) {
		_, err := svc.GetMonitor("non-existent")
		if err == nil {
			t.Fatal("ожидалась ошибка для несуществующего монитора")
		}
	})
}

func TestMonitorService_GetMonitors(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	svc := service.NewMonitorService(repo)

	_ = repo.CreateMonitor(&model.Monitor{ID: "1", URL: "https://one.com"})
	_ = repo.CreateMonitor(&model.Monitor{ID: "2", URL: "https://two.com"})

	monitors, err := svc.GetMonitors()
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if len(monitors) != 2 {
		t.Errorf("ожидалось 2 монитора, получено: %d", len(monitors))
	}
}

func TestMonitorService_DeleteMonitor(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	svc := service.NewMonitorService(repo)

	monitorID := "to-delete"
	_ = repo.CreateMonitor(&model.Monitor{ID: monitorID, URL: "https://test.com"})

	err := svc.DeleteMonitor(monitorID)
	if err != nil {
		t.Fatalf("неожиданная ошибка при удалении: %v", err)
	}

	_, err = svc.GetMonitor(monitorID)
	if err == nil {
		t.Error("ожидалась ошибка после удаления монитора")
	}
}

func TestMonitorService_CreateCheck(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	svc := service.NewMonitorService(repo)

	monitorID := "monitor-for-check"
	_ = repo.CreateMonitor(&model.Monitor{ID: monitorID, URL: "https://test.com"})

	t.Run("успешный чек (2xx) и передача в repository", func(t *testing.T) {
		before := time.Now()

		check, err := svc.CreateCheck(monitorID, 200, 150*time.Millisecond)

		after := time.Now()
		if err != nil {
			t.Fatalf("неожиданная ошибка: %v", err)
		}

		if !check.Success {
			t.Error("ожидалось Success == true для статуса 200")
		}
		if check.StatusCode != 200 {
			t.Errorf("неверный StatusCode: got %d", check.StatusCode)
		}
		if check.Latency != 150*time.Millisecond {
			t.Errorf("неверный Latency: got %v", check.Latency)
		}
		if check.MonitorID != monitorID {
			t.Errorf("неверный MonitorID: got %q", check.MonitorID)
		}
		if check.Time.Before(before) || check.Time.After(after) {
			t.Error("время Time должно быть текущим")
		}

		// Проверка, что Check действительно передан и сохранен в repository
		history, err := repo.GetHistory(monitorID)
		if err != nil {
			t.Fatalf("ошибка получения истории из репозитория: %v", err)
		}
		if len(history) != 1 {
			t.Fatalf("ожидалось 1 запись в репозитории, получено: %d", len(history))
		}
		if history[0].StatusCode != 200 {
			t.Error("репозиторий не получил корректный объект Check")
		}
	})

	t.Run("неуспешный чек (не 2xx)", func(t *testing.T) {
		check, err := svc.CreateCheck(monitorID, 404, 50*time.Millisecond)
		if err != nil {
			t.Fatalf("неожиданная ошибка: %v", err)
		}
		if check.Success {
			t.Error("ожидалось Success == false для статуса 404")
		}
		if check.StatusCode != 404 {
			t.Errorf("неверный StatusCode: got %d", check.StatusCode)
		}
	})
}

func TestMonitorService_GetHistory(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	svc := service.NewMonitorService(repo)

	monitorID := "monitor-history"
	_ = repo.CreateMonitor(&model.Monitor{ID: monitorID, URL: "https://test.com"})

	_, _ = svc.CreateCheck(monitorID, 200, 10*time.Millisecond)
	_, _ = svc.CreateCheck(monitorID, 500, 20*time.Millisecond)

	history, err := svc.GetHistory(monitorID)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if len(history) != 2 {
		t.Fatalf("ожидалось 2 записи в истории, получено: %d", len(history))
	}
	if history[0].StatusCode != 200 || history[1].StatusCode != 500 {
		t.Errorf("неверный порядок или данные истории: got %+v", history)
	}

	t.Run("ошибка при несуществующем мониторе", func(t *testing.T) {
		_, err := svc.GetHistory("non-existent-monitor")
		if err == nil {
			t.Fatal("ожидалась ошибка для несуществующего монитора")
		}
	})
}
