package main

import (
	"net/http"
)

func (app *App) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", app.Ping)
	mux.HandleFunc("/register", app.Register)
	return mux
}
