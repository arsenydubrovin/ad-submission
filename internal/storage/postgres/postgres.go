package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	DB *sql.DB
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
		serial_number INT NOT NULL,
	  advert_id INT REFERENCES adverts(id) ON DELETE CASCADE);
  `)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}

	return &Storage{DB: db}, nil
}

func (s *Storage) CreateAdvert(title, description string, photoLinks *[]string, price int) (err error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	// Commit or rollback depending on result of transaction
	// Named return value (err) is required
	defer func() {
		if err == nil {
			err = tx.Commit()
		}
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var advertId int

	err = tx.QueryRow("INSERT INTO adverts (title, description, price) VALUES ($1, $2, $3) RETURNING id",
		title, description, price).Scan(&advertId)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO photo_links (link, serial_number, advert_id) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i, link := range *photoLinks {
		_, err = stmt.Exec(link, i+1, advertId)
		if err != nil {
			return err
		}
	}

	return nil
}
