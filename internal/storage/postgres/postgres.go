package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/arsenydubrovin/ad-submission/internal/storage"
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
		price INT NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW());
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

func (s *Storage) CreateAdvert(title, description string, photoLinks *[]string, price int) (id int, err error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return 0, err
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
		return 0, err
	}

	stmt, err := tx.Prepare("INSERT INTO photo_links (link, serial_number, advert_id) VALUES ($1, $2, $3)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	for i, link := range *photoLinks {
		_, err = stmt.Exec(link, i+1, advertId)
		if err != nil {
			return 0, err
		}
	}

	return advertId, nil
}

func (s *Storage) GetAdvert(id int) (*storage.Advert, error) {
	if id < 1 {
		return nil, storage.ErrRecordNotFound
	}

	var advert storage.Advert

	stmt := `SELECT id, title, description, price, created_at
					 FROM adverts
					 WHERE id = $1`

	err := s.DB.QueryRow(stmt, id).Scan(
		&advert.Id,
		&advert.Title,
		&advert.Description,
		&advert.Price,
		&advert.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, storage.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	var photoLinks []string

	stmt = `SELECT link
				  FROM photo_links
				  WHERE advert_id = $1
				  ORDER BY serial_number`

	rows, err := s.DB.Query(stmt, advert.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var link string
		err := rows.Scan(&link)
		if err != nil {
			return nil, err
		}
		photoLinks = append(photoLinks, link)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	advert.PhotoLinks = photoLinks

	return &advert, nil
}
