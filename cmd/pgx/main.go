package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

func main() {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	conn, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@127.0.0.1:5432/%s?sslmode=disable", user, password, user))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	var shape []string
	err = conn.QueryRow(context.Background(), "select array_agg(distinct shape) from sighting where shape is not null").Scan(&shape)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%#v\n", shape)
}
