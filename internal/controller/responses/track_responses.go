package responses

import (
	"netradio/internal/model"
	"time"
)

type GetTrackResponse struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Performancer string        `json:"performancer"`
	Year         int           `json:"year"`
	Duration     time.Duration `json:"duration"`
	Audio        string        `json:"audio"`
	Liked        bool          `json:"liked"`
	Found        bool          `json:"-"`
}

type GetTrackListResponse struct {
	Tracks []model.Track `json:"tracks"`
	Count  int           `json:"count"`
}

type CreateTrackResponse struct {
	ID string `json:"id"`
}

type UpdateTrackResponse struct {
	Found bool `json:"-"`
}

type GetTrackCommentsResponse struct {
	Comments []model.Comment `json:"comments"`
}

type CommentTrackResponse struct {
	ID int `json:"id"`
}
