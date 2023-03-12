package model

import "time"

type News struct {
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	PublicationDate time.Time `json:"publication_date"`
	Content         string    `json:"content"`
}

type NewsShortInfo struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	PublicationDate time.Time `json:"date"`
}
