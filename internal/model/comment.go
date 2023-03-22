package model

import "time"

type Comment struct {
	ID       int       `json:"id"`
	UserID   int       `json:"userId"`
	Text     string    `json:"text"`
	Date     time.Time `json:"date"`
	Parent   int       `json:"parent"`
	Children []Comment `json:"children"`
}
