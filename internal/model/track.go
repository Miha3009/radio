package model

import "time"

type Track struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Performancer string        `json:"performancer"`
	Year         int           `json:"year"`
	Audio        string        `json:"audio"`
	Duration     time.Duration `json:"duration"`
}
