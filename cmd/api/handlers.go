package main

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/nelthaarion/go-auth-sample/cmd/api/utils"
	"golang.org/x/crypto/bcrypt"
)

type RequestPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type JWTClaim struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
type LoginResponse struct {
	User  any    `json:"user"`
	Token string `json:"token"`
}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (app *App) Ping(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(os.Getenv("DSN")))
}

func (app *App) Register(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
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
		utils.MessageResponse(w, err.Error(), 400)
	} else {
		utils.MessageResponse(w, "user added successfully", 200)
	}
}

func (app *App) Login(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if req.Method != http.MethodPost {
		utils.MethodNotAllowed(w)
		return
	}

	var loginReq LoginRequest
	if err := json.NewDecoder(req.Body).Decode(&loginReq); err != nil {
		utils.MessageResponse(w, "bad input data", 400)
		return
	}
	user := app.Data.User
	user.Username = loginReq.Username

	foundUser, err := user.GetUser()
	if err != nil {
		utils.MessageResponse(w, "user not found", 404)
		return
	}

	errChan := make(chan error)
	go func(fup, up string) {
		errChan <- bcrypt.CompareHashAndPassword([]byte(fup), []byte(up))
	}(foundUser.Password, loginReq.Password)
	err = <-errChan
	close(errChan)

	if err != nil {
		utils.MessageResponse(w, "passowrd is wrong", 400)
		return
	}
	tokenChan := make(chan string, 1)
	defer close(tokenChan)
	go genereteToken(foundUser.Username, tokenChan)
	token := <-tokenChan
	log.Println(token)
	if len(token) == 0 {
		utils.MessageResponse(w, "somthing went wrong", 500)
		return
	}
	utils.JsonResponse(w, "login successfull", &LoginResponse{User: foundUser, Token: token}, 200)
}

func (app *App) CheckAuthorization(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "application/json")
	if req.Method != http.MethodGet {
		utils.MethodNotAllowed(w)
		return
	}
	authHeader := req.Header.Get("authorization")
	if len(authHeader) == 0 {
		utils.MessageResponse(w, "UNAUTHORIZED", 401)
	}
	splitedHeader := strings.Split(authHeader, " ")
	if len(splitedHeader) < 2 {
		utils.MessageResponse(w, "UNAUTHORIZED", 401)
	}
	ok, err := validateToken(splitedHeader[1])
	if !ok || err != nil {
		utils.MessageResponse(w, err.Error(), 401)
		return
	}
	utils.MessageResponse(w, "AUTHORIZED", 200)
}

func genereteToken(username string, ch chan string) {
	expireAt := time.Now().Add(time.Hour * 12)
	claims := &JWTClaim{Username: username, StandardClaims: jwt.StandardClaims{ExpiresAt: expireAt.Unix()}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT")))

	if err != nil {
		log.Println(err, []byte("STRONGPASSWORD"))
		ch <- ""
	} else {
		ch <- tokenString
	}
}

func validateToken(tokenString string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT")), nil
	})
	if err != nil {
		return false, err
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return false, errors.New("parsing claims failed")
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {

		return false, errors.New("token expired")
	}
	return true, nil
}
