package utils

import (
	"encoding/json"
	"net/http"
)

func MethodNotAllowed(w http.ResponseWriter) {

	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Header().Set("Content-Type", "application/json")
	msg := map[string]string{}
	msg["message"] = "Method not allowed"
	jsonData, _ := json.Marshal(msg)
	w.Write(jsonData)

}

func BadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	msg := map[string]string{}
	msg["message"] = "Bad inputs data"
	jsonData, _ := json.Marshal(msg)
	w.Write(jsonData)
}
