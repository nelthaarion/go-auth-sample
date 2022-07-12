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

func MessageResponse(w http.ResponseWriter, msg string, statusCode int) {
	type Response struct {
		Message string `json:"message"`
	}
	resp := Response{Message: msg}
	jsonData, _ := json.Marshal(resp)
	w.WriteHeader(statusCode)
	w.Write(jsonData)
}

func JsonResponse(w http.ResponseWriter, msg string, data any, statusCode int) {
	type Response struct {
		Message string `json:"message"`
		Data    any    `json:"data"`
	}
	resp := Response{Message: msg, Data: data}
	jsonData, _ := json.Marshal(resp)
	w.WriteHeader(statusCode)
	w.Write(jsonData)
}
