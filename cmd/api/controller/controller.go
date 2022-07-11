package controller

import "net/http"

func Ping(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("PONG"))
}
