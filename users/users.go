package users

import (
	"context"

	"github.com/jackc/pgx/v4"
)

// User is a user
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// UserDirectory allows you interact with users in the database.
type UserDirectory struct {
	conn *pgx.Conn
}

// FindByID attempts to load a user from the database for the provided ID.
func (ud *UserDirectory) FindByID(id int64) (*User, error) {
	var name string
	err := ud.conn.QueryRow(context.Background(), "select name from users where id=$1", id).Scan(&name)
	if err != nil {
		return nil, err
	}

	return &User{id, name}, nil
}

// New returns a new users directory
func New(conn *pgx.Conn) *UserDirectory {
	return &UserDirectory{conn}
}
