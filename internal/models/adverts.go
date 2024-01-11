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

	stmt := `SELECT id, title, description, price, photo_links, created_at
					 FROM adverts
					 WHERE id = $1`

	var advert Advert

	err := am.DB.QueryRow(stmt, id).Scan(
		&advert.Id,
		&advert.Title,
		&advert.Description,
		&advert.Price,
		pq.Array(&advert.PhotoLinks),
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

	return &advert, nil
}

func (am *AdvertModel) GetAll() ([]*Advert, error) {
	stmt := `SELECT id, title, description, price, photo_links, created_at
					 FROM adverts
					 ORDER BY id`

	rows, err := am.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	adverts := []*Advert{}

	for rows.Next() {
		var advert Advert

		err := rows.Scan(
			&advert.Id,
			&advert.Title,
			&advert.Description,
			&advert.Price,
			pq.Array(&advert.PhotoLinks),
			&advert.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		adverts = append(adverts, &advert)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return adverts, nil
}
