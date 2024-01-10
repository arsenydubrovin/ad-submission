package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func openDB(host, port, user, dbName string) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", host, port, user, dbName)

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
	// Create adverts table
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS adverts (
		id SERIAL PRIMARY KEY,
		title VARCHAR(200) NOT NULL,
		description VARCHAR(1000),
		price INT NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW());
  `)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	// Create photo_links table
	stmt, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS photo_links (
		id SERIAL PRIMARY KEY,
		link VARCHAR(255) NOT NULL,
		serial_number INT NOT NULL,
	  advert_id INT REFERENCES adverts(id) ON DELETE CASCADE);
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
