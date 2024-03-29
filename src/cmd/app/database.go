package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = prepareDB(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func prepareDB(db *sql.DB) error {
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS adverts (
		id serial PRIMARY KEY,
		title text NOT NULL,
		description text,
		price int NOT NULL,
		photo_links text[] NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW());
  `)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}
