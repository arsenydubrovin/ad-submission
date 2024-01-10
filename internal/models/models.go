package models

import (
	"database/sql"
	"errors"
)

var ErrRecordNotFound = errors.New("record not found")

type Models struct {
	Adverts AdvertModel
}

func New(db *sql.DB) Models {
	return Models{
		Adverts: AdvertModel{DB: db},
	}
}
