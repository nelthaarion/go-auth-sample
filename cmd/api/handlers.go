package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/nelthaarion/go-auth-sample/cmd/api/utils"
)

type RequestPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (app *App) Ping(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(os.Getenv("DSN")))
}

func (app *App) Register(w http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		utils.MethodNotAllowed(w)
		return
	}

	user := app.Data.User
	payload := RequestPayload{}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		utils.BadRequest(w)
		return
	}

	user.Username = payload.Username
	user.Password = payload.Password
	if err := user.CreateNewUser(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write([]byte("successfull"))

	}
}
