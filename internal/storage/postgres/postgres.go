package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(host, port, user, dbName string) (*Storage, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", host, port, user, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Create adverts table
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS adverts (
		id SERIAL PRIMARY KEY,
		title VARCHAR(200) NOT NULL,
		description VARCHAR(1000),
		price INT NOT NULL);
	`)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}

	// Create photo_links table
	stmt, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS photo_links (
		id SERIAL PRIMARY KEY,
		link VARCHAR(255) NOT NULL,
		serial_number INT NOT NULL);
  `)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}
