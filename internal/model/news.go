package model

type News struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	PublicationTime int64  `json:"publication_date"`
	Content         string `json:"content"`
}
