package middleware

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"gitlab.com/leadcycl/confluence/users"
)

type response struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	StatusText string `json:"statusText"`
}

// Middleware stores the database connection
type Middleware struct {
	conn *pgx.Conn
}

// New returns an instance of Middleware
func New(conn *pgx.Conn) (*Middleware, error) {
	return &Middleware{
		conn,
	}, nil
}

// Handler checks if the provided user id exists
func (m *Middleware) Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		jsonResponse(w, r, "Not Found", http.StatusBadRequest)
		return
	}

	// logrus.Infof("New Request:\nuser agent: %v\nremote address: %v\nmethod: %v\n\n", r.UserAgent(), r.RemoteAddr, r.Method)
	// logrus.Infof("New Request, path: %v query: %v", r.URL.Path, r.URL.Query())

	queryID := r.URL.Query().Get("id")
	if queryID == "" {
		logrus.Errorf("missing queryID")
		jsonResponse(w, r, "Not Found", http.StatusBadRequest)
		return
	}

	n, err := strconv.ParseInt(queryID, 10, 64)
	if err != nil {
		logrus.Errorf("unable to convert query.id to int: %v", err)
		jsonResponse(w, r, "Not Found", http.StatusBadRequest)
		return
	}

	// logrus.Infof("queryId: %v", queryID)
	// u, err := *users.FindByID(n)

	ud := users.New(m.conn)
	u, err := ud.FindByID(n)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		jsonResponse(w, r, "Not Found", http.StatusNotFound)
		return
	}

	logrus.Infof("found user: %v", *u)

	jsonResponse(w, r, "Success", http.StatusOK)
}

// FooHandler is just for testing
func (m *Middleware) FooHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello foo!, %q", html.EscapeString(r.URL.Path))
}

// CatchAllHandler will return a not found error
func (m *Middleware) CatchAllHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CatchAllHandler caught path:", r.URL.Path)
	jsonResponse(w, r, "not found", 400)
}

func jsonResponse(w http.ResponseWriter, r *http.Request, message string, statusCode int) {
	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
	}

	// jsonUser, err := json.Marshal(&user{id, name})
	jsonRes, err := json.Marshal(&response{
		Message:    message,
		StatusText: http.StatusText(statusCode),
		StatusCode: statusCode,
	})
	if err != nil {
		logrus.Errorf("unable to marshal user json: %v", err)
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonRes))
}
