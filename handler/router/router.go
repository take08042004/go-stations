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

	protected := middleware.Recovery(
		middleware.UserAgentMiddleWare(
			middleware.LoggingMiddleware(
				middleware.BasicAuthMiddleware(todoHandler),
			),
		),
	)

	mux.Handle("/todos", protected)

	panicHandler := middleware.PanicHandler{}
	safeHandler := middleware.Recovery(panicHandler)

	mux.Handle("/do-panic",safeHandler)

	handler := http.HandlerFunc(SampleHandler)

	//middlewareチェーン
	wrapped := middleware.UserAgentMiddleWare(handler)

	mux.Handle("/os", wrapped)


	//middlewareチェーン
	wrapped1 := middleware.UserAgentMiddleWare(handler)
	wrapped2 := middleware.LoggingMiddleware(wrapped1)
	wrapped3 := middleware.Recovery(wrapped2)

	mux.Handle("/test", wrapped3)

	return mux
}

func SampleHandler(w http.ResponseWriter, r *http.Request) {
	os, ok := r.Context().Value(middleware.OSKey).(string)

	if !ok {
		os = "unknown"
	}

	w.Write([]byte("Your OS is " + os))
}
