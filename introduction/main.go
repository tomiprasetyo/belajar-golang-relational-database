package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
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

	db, err := sql.Open("pgx", dsn.String())
	if err != nil {
		fmt.Println("sql.Open", err)
		return
	}

	defer func() {
		_ = db.Close()
		fmt.Println("closed")
	}()

	if err := db.PingContext(context.Background()); err != nil {
		fmt.Println("db.PingContext", err)
		return
	}

	row := db.QueryRowContext(context.Background(), "SELECT birth_year FROM users WHERE name = 'tomi'")
	if err := row.Err(); err != nil {
		fmt.Println("db.QueryContext", err)
		return
	}

	var birth_year int64

	if err := row.Scan(&birth_year); err != nil {
		fmt.Println("row.Scan", err)
		return
	}

	fmt.Println("birth_year", birth_year)

	// -

	// rawan terkena sql injection
	// if _, err := db.ExecContext(context.Background(), "INSERT INTO users(name, birth_year) VALUES('prasetyo', 1889)"); err != nil {
	// 	fmt.Println("db.ExecContext", err)
	// 	return
	// }

	// -

	// menggunakan placeholder '$'
	name := "Tomi Prasetyo"
	birth_year = 1996

	_, err = db.ExecContext(context.Background(), "INSERT INTO users(name, birth_year) VALUES($1, $2)", name, birth_year)

	if err != nil {
		fmt.Println("db.ExecContext", err)
		return
	}

	// -

	rows, err := db.QueryContext(context.Background(), "SELECT name, birth_year FROM users")
	if err != nil {
		fmt.Println("row.Scan", err)
		return
	}

	defer func() {
		_ = rows.Close()
	}()

	if rows.Err() != nil {
		fmt.Println("row.Err()", err)
		return
	}

	for rows.Next() {
		var name string
		var birth_year int64

		if err := rows.Scan(&name, &birth_year); err != nil {
			fmt.Println("rows.Scan", err)
			return
		}

		fmt.Println("name", name, "birth year", birth_year)
	}
}
