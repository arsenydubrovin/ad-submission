package storage

import (
	"errors"
	"time"
)

var ErrRecordNotFound = errors.New("record not found")

type Advert struct {
	Id          int
	Title       string
	Description string
	PhotoLinks  []string
	Price       int
	CreatedAt   time.Time
}
