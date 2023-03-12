package responses

import (
	"netradio/internal/model"
	"time"
)

type GetNewsListResponse struct {
	NewsList []model.NewsShortInfo `json:"news"`
	Count    int                   `json:"count"`
}

type GetNewsResponse struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	PublicationDate time.Time `json:"publication_date"`
	Found           bool      `json:"-"`
}

type UpdateNewsResponse struct {
	Found bool `json:"-"`
}
