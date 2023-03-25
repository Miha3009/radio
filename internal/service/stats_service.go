package service

import (
	"netradio/internal/controller/requests"
	"netradio/internal/controller/responses"
	"netradio/internal/repository"
)

type StatsService interface {
	GetListenersStat(r requests.GetListenersStatRequest) (responses.GetListenersStatResponse, error)
}

func NewStatsService() StatsService {
	return &StatsServiceImpl{
		db: repository.NewStatsDB(),
	}
}

type StatsServiceImpl struct {
	db repository.StatsDB
}

func (s *StatsServiceImpl) GetListenersStat(r requests.GetListenersStatRequest) (responses.GetListenersStatResponse, error) {
	var res responses.GetListenersStatResponse
	stats, err := s.db.GetListeners(r.ChannelID, r.From, r.To)
	if err != nil {
		return res, err
	}
	res.Stats = stats

	return res, nil
}
