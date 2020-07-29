package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"gitlab.com/leadcycl/confluence/middleware"
)

// Router is exported and used in main.go
func Router(conn *pgx.Conn) (*mux.Router, error) {
	m, err := middleware.New(conn)
	if err != nil {
		return nil, err
	}

	router := mux.NewRouter()

	// router.HandleFunc("/api/task", middleware.GetAllTask).Methods("GET", "OPTIONS")
	// router.HandleFunc("/api/task", middleware.CreateTask).Methods("POST", "OPTIONS")
	// router.HandleFunc("/api/task/{id}", middleware.TaskComplete).Methods("PUT", "OPTIONS")
	// router.HandleFunc("/api/undoTask/{id}", middleware.UndoTask).Methods("PUT", "OPTIONS")
	// router.HandleFunc("/api/deleteTask/{id}", middleware.DeleteTask).Methods("DELETE", "OPTIONS")
	// router.HandleFunc("/api/deleteAllTask", middleware.DeleteAllTask).Methods("DELETE", "OPTIONS")

	router.HandleFunc("/", m.Handler).Methods("GET", "OPTIONS")
	router.HandleFunc("/foo", m.FooHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/redir", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/foo", http.StatusMovedPermanently)
	})
	router.PathPrefix("/").HandlerFunc(m.CatchAllHandler)

	return router, nil
}
