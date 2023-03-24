package service

import (
	"netradio/internal/controller/requests"
	"netradio/internal/repository"
)

type ScheduleService interface {
	DeleteTrack(r requests.DeleteTrackFromScheduleRequest) error
	UpdateTracks(r requests.UpdateTracksFromScheduleRequest) error
}

func NewScheduleService() ScheduleService {
	return &ScheduleServiceImpl{
		db: repository.NewScheduleDB(),
	}
}

type ScheduleServiceImpl struct {
	db repository.ScheduleDB
}

func (s *ScheduleServiceImpl) DeleteTrack(r requests.DeleteTrackFromScheduleRequest) error {
	return s.db.DeleteTrack(r.ID)
}

func (s *ScheduleServiceImpl) UpdateTracks(r requests.UpdateTracksFromScheduleRequest) error {
	return s.db.UpdateTracks(r.Tracks)
}
