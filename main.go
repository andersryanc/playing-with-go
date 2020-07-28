package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"gitlab.com/leadcycl/confluence/router"
)

func main() {
	port := flag.Int("port", 8080, "the port to listen on")
	flag.Parse()

	r := router.Router()
	logrus.Infof("Listening @ http://localhost:%d\n", *port)
	logrus.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), r))
}
