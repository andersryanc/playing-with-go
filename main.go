package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"gitlab.com/leadcycl/confluence/users"
)

var conn *pgx.Conn

type response struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	StatusText string `json:"statusText"`
}

func main() {
	port := flag.Int("port", 8080, "the port to listen on")
	flag.Parse()

	var err error
	conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	logrus.Infof("Listening @ http://localhost:%d\n", *port)

	http.HandleFunc("/", handler)
	http.HandleFunc("/foo", fooHandler)
	http.HandleFunc("/redir", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/foo", http.StatusMovedPermanently)
	})
	logrus.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
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

	ud := users.New(conn)
	u, err := ud.FindByID(n)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		jsonResponse(w, r, "Not Found", http.StatusNotFound)
		return
	}

	logrus.Infof("found user: %v", *u)

	jsonResponse(w, r, "Success", http.StatusOK)
}

func fooHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello foo!, %q", html.EscapeString(r.URL.Path))
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
