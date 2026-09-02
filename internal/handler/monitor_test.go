package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/enlivee/pulse/internal/handler"
	"github.com/enlivee/pulse/internal/repository"
	"github.com/enlivee/pulse/internal/service"
)

func setupHandler(t *testing.T) *handler.MonitorHandler {
	t.Helper()
	repo := repository.NewInMemoryRepository()
	svc := service.NewMonitorService(repo)
	return handler.NewMonitorHandler(svc)
}

func TestCreateMonitor(t *testing.T) {
	h := setupHandler(t)

	t.Run("успешное создание", func(t *testing.T) {
		reqBody := struct {
			URL      string        `json:"url"`
			Interval time.Duration `json:"interval"`
		}{
			URL:      "https://example.com",
			Interval: 30 * time.Second,
		}
		b, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/monitors", bytes.NewReader(b))
		w := httptest.NewRecorder()

		h.CreateMonitor(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("ожидался статус %d, получен %d: %s", http.StatusCreated, w.Code, w.Body.String())
		}
	})

	t.Run("невалидный JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/monitors", bytes.NewReader([]byte("{invalid}")))
		w := httptest.NewRecorder()

		h.CreateMonitor(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("ожидался статус %d, получен %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestGetMonitors(t *testing.T) {
	t.Run("список с мониторами", func(t *testing.T) {
		repo := repository.NewInMemoryRepository()
		svc := service.NewMonitorService(repo)
		h := handler.NewMonitorHandler(svc)

		_, _ = svc.CreateMonitor("https://one.com", 10*time.Second)
		_, _ = svc.CreateMonitor("https://two.com", 20*time.Second)

		req := httptest.NewRequest(http.MethodGet, "/monitors", nil)
		w := httptest.NewRecorder()

		h.GetMonitors(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("ожидался статус %d, получен %d", http.StatusOK, w.Code)
		}
		if w.Header().Get("Content-Type") != "application/json" {
			t.Errorf("ожидался Content-Type application/json")
		}

		var monitors []map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &monitors); err != nil {
			t.Fatalf("не удалось распарсить ответ: %v", err)
		}
		if len(monitors) != 2 {
			t.Errorf("ожидалось 2 монитора, получено %d", len(monitors))
		}
	})

	t.Run("пустой список", func(t *testing.T) {
		h := setupHandler(t)
		req := httptest.NewRequest(http.MethodGet, "/monitors", nil)
		w := httptest.NewRecorder()

		h.GetMonitors(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("ожидался статус %d, получен %d", http.StatusOK, w.Code)
		}

		var monitors []map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &monitors); err != nil {
			t.Fatalf("не удалось распарсить ответ: %v", err)
		}
		if len(monitors) != 0 {
			t.Errorf("ожидался пустой список, получено %d элементов", len(monitors))
		}
	})
}

func TestGetMonitor(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	svc := service.NewMonitorService(repo)
	h := handler.NewMonitorHandler(svc)

	created, _ := svc.CreateMonitor("https://test.com", 15*time.Second)

	t.Run("успешное получение", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/monitors/"+created.ID, nil)
		req.SetPathValue("id", created.ID)
		w := httptest.NewRecorder()

		h.GetMonitor(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("ожидался статус %d, получен %d", http.StatusOK, w.Code)
		}
		if w.Header().Get("Content-Type") != "application/json" {
			t.Errorf("ожидался Content-Type application/json")
		}

		var monitor map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &monitor); err != nil {
			t.Fatalf("не удалось распарсить ответ: %v", err)
		}
		if monitor["url"] != "https://test.com" {
			t.Errorf("неверный URL в ответе: got %v", monitor["url"])
		}
	})

	t.Run("монитор не найден", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/monitors/non-existent", nil)
		req.SetPathValue("id", "non-existent")
		w := httptest.NewRecorder()

		h.GetMonitor(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("ожидался статус %d, получен %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestDeleteMonitor(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	svc := service.NewMonitorService(repo)
	h := handler.NewMonitorHandler(svc)

	created, _ := svc.CreateMonitor("https://delete-me.com", 10*time.Second)

	t.Run("успешное удаление", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/monitors/"+created.ID, nil)
		req.SetPathValue("id", created.ID)
		w := httptest.NewRecorder()

		h.DeleteMonitor(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("ожидался статус %d, получен %d", http.StatusNoContent, w.Code)
		}

		_, err := svc.GetMonitor(created.ID)
		if err == nil {
			t.Error("монитор должен был быть удален")
		}
	})
}
