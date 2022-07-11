package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type App struct {
	DB *sql.DB
}

func main() {
	conn := getDBConnection()
	if conn == nil {
		log.Panic("Could not connect to database")
	}
	app := App{DB: conn}

	log.Fatal(http.ListenAndServe(":80", app.Routes()))
}

func getDBConnection() *sql.DB {
	DSN := os.Getenv("DSN")
	counter := 0

	for {

		db, err := sql.Open("pgx", DSN)

		if err != nil {
			log.Println(err)
			counter++
		} else {
			return db
		}

		if counter < 10 {
			time.Sleep(time.Second * 3)
			continue
		} else {
			return nil
		}

	}
}
