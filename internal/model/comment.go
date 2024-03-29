package model

import "time"

type Comment struct {
	ID         int       `json:"id"`
	UserID     int       `json:"userId"`
	UserAvatar string    `json:"userAvatar"`
	UserName   string    `json:"userName"`
	Text       string    `json:"text"`
	Date       time.Time `json:"date"`
	Parent     int       `json:"parent"`
	Children   []Comment `json:"children"`
}
