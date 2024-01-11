package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/arsenydubrovin/ad-submission/internal/validator"
	"github.com/lib/pq"
)

type Advert struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PhotoLinks  []string  `json:"photoLinks"`
	Price       int       `json:"price"`
	CreatedAt   time.Time `json:"-"`
}

func ValidateAdvert(v *validator.Validator, advert *Advert) {
	v.Check(advert.Title != "", "title", "must be provided")
	v.Check(len([]rune(advert.Title)) <= 200, "title", "must be no more than 200 characters")

	// description can be empty
	v.Check(len([]rune(advert.Description)) <= 1000, "description", "must be no more than 1000 characters")

	v.Check(advert.PhotoLinks != nil, "photoLinks", "must be provided")
	for _, link := range advert.PhotoLinks {
		v.Check(link != "", "links", "must be non-empty")
	}
	v.Check(len(advert.PhotoLinks) >= 1, "photoLinks", "must contain at least 1 link")
	v.Check(len(advert.PhotoLinks) <= 3, "photoLinks", "must contain no more than 3 links")

	v.Check(advert.Price > 0, "price", "must be positive")
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
