package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/leadcycl/confluence/middleware"
)

// Router is exported and used in main.go
func Router() *mux.Router {
	router := mux.NewRouter()

	// router.HandleFunc("/api/task", middleware.GetAllTask).Methods("GET", "OPTIONS")
	// router.HandleFunc("/api/task", middleware.CreateTask).Methods("POST", "OPTIONS")
	// router.HandleFunc("/api/task/{id}", middleware.TaskComplete).Methods("PUT", "OPTIONS")
	// router.HandleFunc("/api/undoTask/{id}", middleware.UndoTask).Methods("PUT", "OPTIONS")
	// router.HandleFunc("/api/deleteTask/{id}", middleware.DeleteTask).Methods("DELETE", "OPTIONS")
	// router.HandleFunc("/api/deleteAllTask", middleware.DeleteAllTask).Methods("DELETE", "OPTIONS")

	router.HandleFunc("/", middleware.Handler).Methods("GET", "OPTIONS")
	router.HandleFunc("/foo", middleware.FooHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/redir", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/foo", http.StatusMovedPermanently)
	})

	return router
}
