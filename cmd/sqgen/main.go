package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	sq "github.com/bokwoon95/go-structured-query/postgres"

	"github.com/slcjordan/ufo/tables"
)

func distinctShapes(db *sql.DB) []string {
	var shapes []string
	var curr string

	s := tables.SIGHTING().As("s")
	query := sq.
		SelectDistinct().
		Selectx(
			func(row *sq.Row) { // map row
				curr = row.String(s.SHAPE)
			},
			func() { // agg
				shapes = append(shapes, curr)
			},
		).
		From(s)

	err := query.FetchContext(context.Background(), db)

	if err != nil {
		log.Fatalf("QueryRow failed: %v\n", err)
	}
	return shapes
}

type sightingSummary struct {
	DurationInSeconds int64
	Comments          string
}

func (s sightingSummary) String() string {
	return fmt.Sprintf("(%s) %s", time.Second*time.Duration(s.DurationInSeconds), s.Comments)
}

func utahCTEExample(db *sql.DB) []sightingSummary {
	var result []sightingSummary
	/*
		durationAlias := "cte_duration"
		commentsAlias := "cte_comments"
	*/

	s := tables.SIGHTING().As("s")
	utahSightings := sq.Select(
		s.DURATION_IN_SECONDS, s.COMMENTS,
	).From(s).CTE("utah_sightings")
	query, args := sq.With(utahSightings).
		Select(utahSightings["duration_in_seconds"], utahSightings["comments"]).
		From(utahSightings).Limit(10).ToSQL()

	fmt.Println("\n\n", query)
	rows, err := db.QueryContext(context.Background(), query, args...)
	if err != nil {
		log.Fatalf("QueryRow failed: %v\n", err)
	}
	for rows.Next() {
		var currDuration sql.NullInt64
		var currComment sql.NullString
		err = rows.Scan(&currDuration, &currComment)
		if err != nil {
			log.Fatalf("QueryRow failed: %v\n", err)
		}
		result = append(result, sightingSummary{DurationInSeconds: currDuration.Int64, Comments: currComment.String})
	}
	return result
}

func main() {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	fmt.Println("user is", user)
	fmt.Println("password is", password)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"127.0.0.1", 5432, user, password, user)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Printf("%#v\n", distinctShapes(db))
	for _, summary := range utahCTEExample(db) {
		fmt.Printf("%v\n", summary)
	}
}
