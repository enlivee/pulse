package main

import (
	"net/http"

	"github.com/enlivee/pulse/internal/handler"
	"github.com/enlivee/pulse/internal/repository"
	"github.com/enlivee/pulse/internal/service"
)

func main() {
	repo := repository.NewInMemoryRepository()
	service := service.NewMonitorService(repo)
	handler := handler.NewMonitorHandler(service)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /monitors", handler.GetMonitors)
	mux.HandleFunc("POST /monitors", handler.CreateMonitor)
	mux.HandleFunc("GET /monitors/{id}", handler.GetMonitor)
	mux.HandleFunc("DELETE /monitors/{id}", handler.DeleteMonitor)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
