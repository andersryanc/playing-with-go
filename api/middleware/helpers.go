package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func jsonData(w http.ResponseWriter, r *http.Request, data interface{}) {
	jsonRes, err := json.Marshal(data)
	if err != nil {
		logrus.Errorf("unable to marshal json: %v", err)
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonRes))
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
