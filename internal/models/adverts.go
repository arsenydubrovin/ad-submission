package models

import (
	"database/sql"
	"errors"
	"time"
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

func (am *AdvertModel) Insert(advert *Advert) (id int, err error) {
	tx, err := am.DB.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err == nil {
			err = tx.Commit()
		}
		_ = tx.Rollback()
	}()

	err = tx.QueryRow(
		`INSERT INTO adverts (title, description, price)
		 VALUES ($1, $2, $3)
	   RETURNING id`,
		advert.Title, advert.Description, advert.Price).Scan(&id)
	if err != nil {
		return 0, err
	}

	stmt, err := tx.Prepare(`INSERT INTO photo_links (link, serial_number, advert_id)
													 VALUES ($1, $2, $3)
												 `)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	for i, link := range advert.PhotoLinks {
		_, err = stmt.Exec(link, i+1, id)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
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
