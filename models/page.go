package models

import (
	"time"
)

type Page struct {
	Id int64
	Content string
	Created time.Time
}