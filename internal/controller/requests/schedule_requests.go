package requests

import (
	"netradio/internal/model"
	"time"
)

type AddTrackToScheduleRequest struct {
	ChannelID string
	TrackID   string    `json:"trackid"`
	StartDate time.Time `json:"startdate"`
	EndDate   time.Time `json:"enddate"`
}

type DeleteTrackFromScheduleRequest struct {
	ID string
}

type UpdateTracksFromScheduleRequest struct {
	Tracks []model.ScheduleTrack `json:"tracks"`
}
