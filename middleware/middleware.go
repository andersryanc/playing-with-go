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
		jsonMessage(w, r, "Not Found", http.StatusBadRequest)
		return
	}

	// logrus.Infof("New Request:\nuser agent: %v\nremote address: %v\nmethod: %v\n\n", r.UserAgent(), r.RemoteAddr, r.Method)
	// logrus.Infof("New Request, path: %v query: %v", r.URL.Path, r.URL.Query())

	queryID := r.URL.Query().Get("id")
	if queryID == "" {
		logrus.Errorf("missing queryID")
		jsonMessage(w, r, "Not Found", http.StatusBadRequest)
		return
	}

	n, err := strconv.ParseInt(queryID, 10, 64)
	if err != nil {
		logrus.Errorf("unable to convert query.id to int: %v", err)
		jsonMessage(w, r, "Not Found", http.StatusBadRequest)
		return
	}

	// logrus.Infof("queryId: %v", queryID)
	// u, err := *users.FindByID(n)

	ud := users.New(m.conn)
	u, err := ud.FindByID(n)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		jsonMessage(w, r, "Not Found", http.StatusNotFound)
		return
	}

	logrus.Infof("found user: %v", *u)

	jsonData(w, r, u)
}

// GetAllUsersHandler loads all the users from the database
func (m *Middleware) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	ud := users.New(m.conn)
	data, err := ud.GetAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		jsonMessage(w, r, "Not Found", http.StatusNotFound)
		return
	}

	logrus.Infof("found users: %v", *data)

	jsonData(w, r, data)
}

// FooHandler is just for testing
func (m *Middleware) FooHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello foo!, %q", html.EscapeString(r.URL.Path))
}

// CatchAllHandler will return a not found error
func (m *Middleware) CatchAllHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CatchAllHandler caught path:", r.URL.Path)
	jsonMessage(w, r, "not found", 400)
}

func jsonMessage(w http.ResponseWriter, r *http.Request, message string, statusCode int) {
	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
	}

	// jsonUser, err := json.Marshal(&user{id, name})
	jsonRes, err := json.Marshal(&response{
		Message:    message,
		StatusCode: statusCode,
		StatusText: http.StatusText(statusCode),
	})
	if err != nil {
		logrus.Errorf("unable to marshal response json: %v", err)
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonRes))
}

func jsonData(w http.ResponseWriter, r *http.Request, data interface{}) {
	jsonRes, err := json.Marshal(data)
	if err != nil {
		logrus.Errorf("unable to marshal json: %v", err)
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonRes))
}
