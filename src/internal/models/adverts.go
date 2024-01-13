package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/arsenydubrovin/ad-submission/src/internal/validator"
	"github.com/lib/pq"
)

type Advert struct {
	Id               int       `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description,omitempty"`
	PrimaryPhotoLink string    `json:"primaryPhotoLink"`
	PhotoLinks       []string  `json:"photoLinks,omitempty"`
	Price            int       `json:"price"`
	CreatedAt        time.Time `json:"-"`
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

func (am *AdvertModel) Get(id int, filters Filters) (*Advert, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	stmt := `SELECT id, title, price, created_at, description, photo_links
					 FROM adverts
					 WHERE id = $1`

	var advert Advert

	// Optional fields
	var photoLinks []string
	var description string

	err := am.DB.QueryRow(stmt, id).Scan(
		&advert.Id,
		&advert.Title,
		&advert.Price,
		&advert.CreatedAt,
		&description,
		pq.Array(&photoLinks),
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	advert.PrimaryPhotoLink = photoLinks[0]

	if filters.includeDescription() {
		advert.Description = description
	}

	if filters.includeAllPhotos() {
		advert.PhotoLinks = photoLinks
	}

	return &advert, nil
}

func (am *AdvertModel) GetAll(filters Filters) ([]*Advert, Info, error) {
	stmt := fmt.Sprintf(`SELECT id, title, price, photo_links, created_at as date, count(*) OVER()
					 FROM adverts
					 ORDER BY %s %s, id ASC
					 LIMIT $1 OFFSET $2`,
		filters.sortColumn(),
		filters.sortDirection())

	fmt.Println("limit", filters.limit(), "offset", filters.offset())
	rows, err := am.DB.Query(stmt, filters.limit(), filters.offset())
	if err != nil {
		return nil, Info{}, err
	}
	defer rows.Close()

	var adverts []*Advert
	var totalRecords int

	for rows.Next() {
		var advert Advert
		var photoLinks []string

		err := rows.Scan(
			&advert.Id,
			&advert.Title,
			&advert.Price,
			pq.Array(&photoLinks),
			&advert.CreatedAt,
			&totalRecords,
		)
		if err != nil {
			return nil, Info{}, err
		}

		advert.PrimaryPhotoLink = photoLinks[0]

		adverts = append(adverts, &advert)
	}

	if err = rows.Err(); err != nil {
		return nil, Info{}, err
	}

	info := calculateInfo(filters.Page, filters.PageSize, totalRecords)

	return adverts, info, nil
}
