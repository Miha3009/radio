package responses

import (
	"netradio/internal/model"
	"time"
)

type GetChannelsResponse struct {
	Channels []model.ChannelShortInfo `json:"channels"`
}

type GetChannelResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	Logo        string `json:"logo"`
	Found       bool   `json:"-"`
}

type UpdateChannelResponse struct {
	Found bool `json:"-"`
}

type GetCurrentTrackResponse struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Performancer string        `json:"performancer"`
	Year         int           `json:"year"`
	Duration     time.Duration `json:"duration"`
	Liked        bool          `json:"liked"`
	LikeCount    int           `json:"likeCount"`
	CurrentTime  time.Duration `json:"currentTime"`
}

type GetScheduleResponse struct {
	Past []model.ScheduleTrackFull `json:"past"`
	Next []model.ScheduleTrackFull `json:"next"`
}

type GetScheduleRangeResponse struct {
	Tracks []model.ScheduleTrackFull `json:"tracks"`
}

type UploadLogoResponse struct {
	Logo string `json:"logo"`
}
