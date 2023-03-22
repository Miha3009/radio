package responses

import (
	"netradio/internal/model"
	"time"
)

type GetTrackResponse struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Perfomancer string        `json:"perfomancer"`
	Year        int           `json:"year"`
	Duration    time.Duration `json:"duration"`
	Liked       bool          `json:"liked"`
	Found       bool          `json:"-"`
}

type GetTrackListResponse struct {
	Tracks []model.Track `json:"tracks"`
	Count  int           `json:"count"`
}

type UpdateTrackResponse struct {
	Found bool `json:"-"`
}

type GetTrackCommentsResponse struct {
	Comments []model.Comment `json:"comments"`
}
