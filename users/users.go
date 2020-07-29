package users

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
)

// User is a user
type User struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"createdAt"`
}

// UserDirectory allows you interact with users in the database.
type UserDirectory struct {
	conn *pgx.Conn
}

// New returns an instance of UserDirectory.
func New(conn *pgx.Conn) *UserDirectory {
	return &UserDirectory{conn}
}

// FindByID attempts to load a user from the database for the provided ID.
func (ud *UserDirectory) FindByID(id int64) (*User, error) {
	u := User{}
	u.ID = id
	err := ud.conn.QueryRow(context.Background(), "select name, created_at from users where id=$1", id).Scan(&u.Name, &u.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// GetAll attempts to load all the users from the database.
func (ud *UserDirectory) GetAll() (*[]User, error) {
	var users []User

	rows, err := ud.conn.Query(context.Background(), "select id, name, created_at from users")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		u := User{}
		err = rows.Scan(&u.ID, &u.Name, &u.CreatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &users, nil
}
