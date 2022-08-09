package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v4"
)

type User struct {
	Name      *string
	BirthYear *int64
}

func main() {
	conn, err := newConn()
	if err != nil {
		fmt.Println("newConn", err)
		return
	}

	defer func() {
		_ = conn.Close(context.Background())
		fmt.Println("closed")
	}()

	newName := func(s string) *string {
		return &s
	}

	newBirthYear := func(b int64) *int64 {
		return &b
	}

	if err := insertUsers(conn, []User{
		{
			Name:      newName("Dika"),
			BirthYear: newBirthYear(1989),
		},
		{
			Name:      newName("Tomi"),
			BirthYear: newBirthYear(1800),
		},
	}); err != nil {
		fmt.Println("insertUsers", err)
	}

}

func newConn() (*pgx.Conn, error) {
	dsn := url.URL{
		Scheme: "postgres",
		Host:   "localhost:5432",
		User:   url.UserPassword("tomi", "root"),
		Path:   "tomi",
	}

	// ssl mode
	q := dsn.Query()
	q.Add("sslmode", "disable")

	dsn.RawQuery = q.Encode()

	conn, err := pgx.Connect(context.Background(), dsn.String())
	if err != nil {
		return nil, fmt.Errorf("pgx.Connect %w", err)
	}

	return conn, nil
}

func insertUsers(conn *pgx.Conn, users []User) error {
	if err := conn.BeginFunc(context.Background(), func(tx pgx.Tx) error {
		for _, user := range users {
			_, err := conn.Exec(context.Background(), "INSERT INTO users(name, birth_year) VALUES($1, $2)", user.Name, user.BirthYear)
			if err != nil {
				return fmt.Errorf("tx.ExecContext %w", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("conn.BeginFunc %w", err)
	}

	return nil
}
