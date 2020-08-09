package middleware

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/sirupsen/logrus"
)

func setup() (*Middleware, error) {
	dbURL, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		panic("missing environment var DATABASE_URL")
	}

	dbConfig, err := pgx.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("unable to parse db config: %v", err)
	}
	dbConfig.Logger = logrusadapter.NewLogger(logrus.New())

	conn, err := pgx.ConnectConfig(context.Background(), dbConfig)
	if err != nil {
		// NOTE: error logging is being handled by dbConfig.Logger, so no need to log again here.
		os.Exit(1)
	}

	m, err := New(conn)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func TestHandler(t *testing.T) {
	m, err := setup()
	if err != nil {
		t.Fatalf("unable to create middleware: %v", err)
	}

	req, err := http.NewRequest("GET", "localhost:8080", nil)
	if err != nil {
		t.Fatalf("unable to create request: %v", err)
	}
	rec := httptest.NewRecorder()
	m.Handler(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}

	var data struct {
		Message    string `json:"message"`
		statusCode int
		statusText string
	}
	json.Unmarshal(b, &data)

	// fmt.Printf("got: %v", string(b))
	want := "Hello, world!"
	got := data.Message
	if want != got {
		t.Fatalf("wanted data.Message: %q, got: %q", want, got)
	}
}
