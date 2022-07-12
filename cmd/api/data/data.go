package data

import (
	"context"
	"database/sql"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

const dbTimeout = time.Second * 3

type Model struct {
	User User
}

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Firstname string    `json:"firstname,omitempty"`
	Lastname  string    `json:"lastname,omitempty"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func New(ndb *sql.DB) Model {
	db = ndb
	return Model{User: User{}}
}

func (u *User) CreateNewUser() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	pchan := make(chan []byte, 1)
	go func() {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
		pchan <- hashed
	}()
	_, err := db.QueryContext(ctx, "INSERT INTO users (username , password, created_at, firstname, lastname) values ($1 ,$2 ,$3 ,'' , '')", u.Username, <-pchan, time.Now())
	log.Println(err)
	return err
}

func (u *User) GetUser() (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var foundUser User
	row, err := db.QueryContext(ctx, "SELECT * FROM users where username = $1 LIMIT 1", u.Username)
	if err != nil {
		log.Println(err)
		return foundUser, err
	}
	if row.Next() {
		err = row.Scan(
			&foundUser.ID,
			&foundUser.Username,
			&foundUser.Password,
			&foundUser.Firstname,
			&foundUser.Lastname,
			&foundUser.CreatedAt,
		)
		log.Println(err)
	}
	return foundUser, nil
}
