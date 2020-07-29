package router

import (
	"net/http"

	"github.com/andersryanc/playing-with-go/api/middleware"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

// Router creates a new mux.Router to be used by the server.
func Router(conn *pgx.Conn) (*mux.Router, error) {
	m, err := middleware.New(conn)
	if err != nil {
		return nil, err
	}

	router := mux.NewRouter()

	router.HandleFunc("/", m.Handler).Methods("GET", "OPTIONS")
	router.HandleFunc("/users/{id:[0-9]+}", m.GetUserByIDHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/users", m.GetAllUsersHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/foo", m.FooHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/redir", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/foo", http.StatusMovedPermanently)
	})
	router.PathPrefix("/").HandlerFunc(m.CatchAllHandler)

	return router, nil
}
