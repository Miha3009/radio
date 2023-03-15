package model

import "time"

type Track struct {
	ID          int           `json:"id"`
	Title       string        `json:"title"`
	Perfomancer string        `json:"perfomancer"`
	Year        int           `json:"year"`
	Audio       string        `json:"audio"`
	Duration    time.Duration `json:"duration"`
}
