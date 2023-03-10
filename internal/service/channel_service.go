package service

import (
	"netradio/internal/controller/requests"
	"netradio/internal/controller/responses"
	"netradio/internal/model"
	"netradio/internal/repository"
	"strconv"
)

type ChannelService interface {
	GetChannels() (responses.GetChannelsResponse, error)
	GetChannel(r requests.GetChannelRequest) (responses.GetChannelResponse, error)
	CreateChannel(r requests.CreateChannelRequest) error
	DeleteChannel(r requests.DeleteChannelRequest) error
	UpdateChannel(r requests.UpdateChannelRequest) (responses.UpdateChannelResponse, error)
}

func NewChannelService() ChannelService {
	return &ChannelServiceImpl{
		db: repository.NewChannelDB(),
	}
}

type ChannelServiceImpl struct {
	db repository.ChannelDB
}

func (s *ChannelServiceImpl) GetChannels() (responses.GetChannelsResponse, error) {
	var res responses.GetChannelsResponse
	channels, err := s.db.GetChannels()
	if err != nil {
		return res, err
	}

	res.Channels = channels

	return res, nil
}

func (s *ChannelServiceImpl) GetChannel(r requests.GetChannelRequest) (responses.GetChannelResponse, error) {
	var res responses.GetChannelResponse
	channel, err := s.db.GetChannelById(r.ID)
	if err != nil {
		return res, err
	}
	if channel == nil {
		res.Found = false
		return res, err
	}

	res.Found = true
	res.ID = strconv.Itoa(channel.ID)
	res.Title = channel.Title
	res.Description = channel.Description

	return res, nil
}

func (s *ChannelServiceImpl) CreateChannel(r requests.CreateChannelRequest) error {
	var channel model.ChannelInfo
	channel.Title = r.Title
	channel.Description = r.Description
	channel.Status = model.StoppedChannel
	return s.db.CreateChannel(channel)
}

func (s *ChannelServiceImpl) DeleteChannel(r requests.DeleteChannelRequest) error {
	return s.db.DeleteChannel(r.ID)
}

func (s *ChannelServiceImpl) UpdateChannel(r requests.UpdateChannelRequest) (responses.UpdateChannelResponse, error) {
	var res responses.UpdateChannelResponse
	channel, err := s.db.GetChannelById(r.ID)
	if err != nil {
		return res, err
	}
	if channel == nil {
		res.Found = false
		return res, nil
	}

	res.Found = true
	if r.Title != nil {
		channel.Title = *r.Title
	}
	if r.Description != nil {
		channel.Description = *r.Description
	}

	err = s.db.UpdateChannel(*channel)
	if err != nil {
		return res, err
	}

	return res, nil
}
