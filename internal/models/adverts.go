package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type Advert struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PhotoLinks  []string  `json:"photolinks"`
	Price       int       `json:"price"`
	CreatedAt   time.Time `json:"-"`
}

type AdvertModel struct {
	DB *sql.DB
}

// Insert() method modifies the Advert structure by writing Id and CreatedAt fields
func (am *AdvertModel) Insert(advert *Advert) error {
	stmt := `INSERT INTO adverts (title, description, price, photo_links)
					 VALUES ($1, $2, $3, $4)
					 RETURNING id, created_at`

	args := []any{
		advert.Title,
		advert.Description,
		advert.Price,
		pq.Array(advert.PhotoLinks),
	}

	return am.DB.QueryRow(stmt, args...).Scan(&advert.Id, &advert.CreatedAt)
}

func (am *AdvertModel) Get(id int) (*Advert, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	var advert Advert

	stmt := `SELECT id, title, description, price, created_at
					 FROM adverts
					 WHERE id = $1`

	err := am.DB.QueryRow(stmt, id).Scan(
		&advert.Id,
		&advert.Title,
		&advert.Description,
		&advert.Price,
		&advert.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	var photoLinks []string

	stmt = `SELECT link
				  FROM photo_links
				  WHERE advert_id = $1
				  ORDER BY serial_number`

	rows, err := am.DB.Query(stmt, advert.Id)
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
