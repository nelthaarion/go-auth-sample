package main

import (
	"net/http"

	"github.com/nelthaarion/go-auth-sample/cmd/api/controller"
)

func (app *App) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", controller.Ping)
	return mux
}
