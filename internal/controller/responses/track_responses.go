package responses

import "time"

type GetTrackResponse struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Perfomancer string        `json:"perfomancer"`
	Year        int           `json:"year"`
	Duration    time.Duration `json:"duration"`
	Liked       bool          `json:"liked"`
	Found       bool          `json:"-"`
}

type UpdateTrackResponse struct {
	Found bool `json:"-"`
}
