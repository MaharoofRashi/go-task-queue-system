package http

import (
	"log"
	"net/http"
	"strings"
)

func SetupRoutes(handler *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.Health)

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.SubmitTask(w, r)
		case http.MethodGet:
			handler.ListTasks(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/cancel") && r.Method == http.MethodPost {
			handler.CancelTask(w, r)
			return
		}

		if r.Method == http.MethodGet {
			handler.GetTask(w, r)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.GetStats(w, r)
	})

	mux.HandleFunc("/workers/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.GetWorkerStatus(w, r)
	})

	return loggingMiddleware(mux)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("📥 %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
