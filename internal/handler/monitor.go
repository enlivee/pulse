package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/enlivee/pulse/internal/service"
)

type MonitorHandler struct {
	service *service.MonitorService
}

func NewMonitorHandler(service *service.MonitorService) *MonitorHandler {
	return &MonitorHandler{
		service: service,
	}
}

func (h *MonitorHandler) CreateMonitor(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL      string        `json:"url"`
		Interval time.Duration `json:"interval"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := h.service.CreateMonitor(req.URL, req.Interval); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *MonitorHandler) GetMonitors(w http.ResponseWriter, r *http.Request) {
	monitors, err := h.service.GetMonitors()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(monitors); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MonitorHandler) GetMonitor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	monitor, err := h.service.GetMonitor(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(monitor); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MonitorHandler) DeleteMonitor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.DeleteMonitor(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
