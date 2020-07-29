package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/andersryanc/playing-with-go/users"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// GetUserByIDHandler will (try to) load a user from the database and return it as json.
func (m *Middleware) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	queryID := vars["id"]
	if queryID == "" {
		logrus.Errorf("missing queryID")
		jsonMessage(w, r, "Not Found", http.StatusBadRequest)
		return
	}

	n, err := strconv.ParseInt(queryID, 10, 64)
	if err != nil {
		logrus.Errorf("unable to convert vars.id to int: %v", err)
		jsonMessage(w, r, "Not Found", http.StatusBadRequest)
		return
	}

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

// GetAllUsersHandler loads all the users from the database and returns them as json.
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
