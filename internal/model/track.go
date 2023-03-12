package model

type Track struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Perfomancer string `json:"perfomancer"`
	Year        int    `json:"year"`
}
