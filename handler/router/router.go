package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/service"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()

	healthzHandler := handler.NewHealthzHandler()
	mux.Handle("/healthz", healthzHandler)

	todoService := service.NewTODOService(todoDB)

	todoHandler := handler.NewTODOHandler(todoService)
	mux.Handle("/todos", todoHandler)

	panicHandler := middleware.PanicHandler{}
	safeHandler := middleware.Recovery(panicHandler)

	mux.Handle("/do-panic",safeHandler)

	return mux
}
