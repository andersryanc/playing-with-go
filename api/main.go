package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/andersryanc/playing-with-go/api/router"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/sirupsen/logrus"
)

func main() {
	wait := flag.Duration("graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	port := flag.Int("port", 8080, "the port to listen on")
	flag.Parse()

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

	r, err := router.Router(conn)
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", *port), // "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	logrus.Infof("Listening @ http://localhost:%d\n", *port)

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), *wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	conn.Close(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
