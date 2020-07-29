package middleware

import (
	"fmt"
	"html"
	"net/http"

	"github.com/jackc/pgx/v4"
)

type response struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	StatusText string `json:"statusText"`
}

// Middleware contains the handler funcs to be used by the router.
type Middleware struct {
	conn *pgx.Conn
}

// New returns an instance of Middleware.
func New(conn *pgx.Conn) (*Middleware, error) {
	return &Middleware{conn}, nil
}

// Handler will send a basic hello world response in json.
func (m *Middleware) Handler(w http.ResponseWriter, r *http.Request) {
	jsonMessage(w, r, "Hello, world!", http.StatusOK)
}

// FooHandler is just for testing.
func (m *Middleware) FooHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello foo!, %q", html.EscapeString(r.URL.Path))
}

// CatchAllHandler will return a not found error in json.
func (m *Middleware) CatchAllHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CatchAllHandler caught path:", r.URL.Path)
	jsonMessage(w, r, "not found", 400)
}
