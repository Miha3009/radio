package requests

import (
	"netradio/internal/model"
)

type AddTrackToScheduleRequest struct {
	ChannelID string
	Tracks    []model.ScheduleTrack `json:"tracks"`
}

type DeleteTrackFromScheduleRequest struct {
	ID string
}

type UpdateTracksFromScheduleRequest struct {
	Tracks []model.ScheduleTrack `json:"tracks"`
}
