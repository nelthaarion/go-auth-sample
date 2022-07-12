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

	"github.com/nelthaarion/go-auth-sample/cmd/api/data"
)

type App struct {
	DB   *sql.DB
	Data data.Model
}

func main() {
	conn := getDBConnection()
	if conn == nil {
		log.Panic("Could not connect to database")
	}
	app := App{DB: conn, Data: data.New(conn)}

	log.Fatal(http.ListenAndServe(":80", app.Routes()))
}

func getDBConnection() *sql.DB {
	DSN := os.Getenv("DSN")
	// DSN := "host=localhost port=5432 user=postgres password=password dbname=users sslmode=disable timezone=UTC connect_timeout=5"
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
