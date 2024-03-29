package service

import (
	"errors"
	"fmt"
	"netradio/internal/controller/requests"
	"netradio/internal/controller/responses"
	"netradio/internal/model"
	"netradio/internal/repository"
	"netradio/pkg/hls"
	"netradio/pkg/stats"
	"strconv"
)

type ChannelService interface {
	GetChannels(r requests.GetChannelsRequest) (responses.GetChannelsResponse, error)
	GetChannel(r requests.GetChannelRequest) (responses.GetChannelResponse, error)
	CreateChannel(r requests.CreateChannelRequest) (responses.CreateChannelResponse, error)
	DeleteChannel(r requests.DeleteChannelRequest) error
	UpdateChannel(r requests.UpdateChannelRequest) (responses.UpdateChannelResponse, error)
	StartChannel(r requests.StartChannelRequest) error
	StopChannel(r requests.StopChannelRequest) error
	UploadLogo(r requests.UploadLogoRequest) (responses.UploadLogoResponse, error)
	AddTrack(r requests.AddTrackToScheduleRequest) error
	GetCurrentTrack(r requests.GetCurrentTrackRequest) (responses.GetCurrentTrackResponse, error)
	GetSchedule(r requests.GetScheduleRequest) (responses.GetScheduleResponse, error)
	GetScheduleRange(r requests.GetScheduleRangeRequest) (responses.GetScheduleRangeResponse, error)
	StartStream(r requests.StartStreamRequest) error
	StopStream(r requests.StopStreamRequest) error
}

func NewChannelService() ChannelService {
	return &ChannelServiceImpl{
		db: repository.NewChannelDB(),
	}
}

type ChannelServiceImpl struct {
	db repository.ChannelDB
}

func (s *ChannelServiceImpl) GetChannels(r requests.GetChannelsRequest) (responses.GetChannelsResponse, error) {
	var res responses.GetChannelsResponse
	channels, channelsCount, err := s.db.GetChannels(r.Offset, r.Limit, r.Query, r.Status)
	if err != nil {
		return res, err
	}

	res.Channels = channels
	res.Count = channelsCount

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
	res.Status = int(channel.Status)
	res.Logo = channel.Logo

	return res, nil
}

func (s *ChannelServiceImpl) CreateChannel(r requests.CreateChannelRequest) (responses.CreateChannelResponse, error) {
	var res responses.CreateChannelResponse
	var channel model.ChannelInfo
	channel.Title = r.Title
	channel.Description = r.Description
	channel.Status = model.StoppedChannel

	id, err := s.db.CreateChannel(channel)
	if err != nil {
		return res, err
	}
	res.ID = strconv.Itoa(id)

	hls.StartChannel(res.ID)
	stats.SetChannelStatus(res.ID, channel.Status)
	stats.RunForChannel(res.ID)

	return res, nil
}

func (s *ChannelServiceImpl) DeleteChannel(r requests.DeleteChannelRequest) error {
	stats.SetChannelStatus(r.ID, model.DeletedChannel)
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

func (s *ChannelServiceImpl) StartChannel(r requests.StartChannelRequest) error {
	stats.SetChannelStatus(r.ID, model.ActiveChannel)
	return s.db.ChangeChannelStatus(r.ID, model.ActiveChannel)
}

func (s *ChannelServiceImpl) StopChannel(r requests.StopChannelRequest) error {
	stats.SetChannelStatus(r.ID, model.StoppedChannel)
	return s.db.ChangeChannelStatus(r.ID, model.StoppedChannel)
}

func (s *ChannelServiceImpl) UploadLogo(r requests.UploadLogoRequest) (responses.UploadLogoResponse, error) {
	var res responses.UploadLogoResponse
	err := s.db.ChangeLogo(r.ID, r.Logo)
	if err != nil {
		return res, err
	}

	res.Logo = r.Logo

	return res, nil
}

func (s *ChannelServiceImpl) AddTrack(r requests.AddTrackToScheduleRequest) error {
	for i := range r.Tracks {
		r.Tracks[i].ChannelId = r.ChannelID
	}
	return repository.NewScheduleDB().AddTracksToSchedule(r.Tracks)
}

func (s *ChannelServiceImpl) GetCurrentTrack(r requests.GetCurrentTrackRequest) (responses.GetCurrentTrackResponse, error) {
	var res responses.GetCurrentTrackResponse
	track, err := s.db.GetCurrentTrack(r.ID)
	if err != nil {
		return res, err
	}
	if track == nil {
		return res, errors.New("Track not found")
	}

	res.ID = track.ID
	res.Title = track.Title
	res.Performancer = track.Performancer
	res.Year = track.Year
	res.Duration = track.Duration
	res.CurrentTime = hls.GetCurrentTrackTime(r.ID)
	res.Listeners = stats.GetListeners(r.ID)
	likeCount, err := repository.NewTrackDB().LikeCount(res.ID)
	if err != nil {
		return res, err
	}
	res.LikeCount = likeCount
	if r.UserID != nil {
		liked, err := repository.NewTrackDB().IsTrackLiked(r.ID, *r.UserID)
		if err != nil {
			return res, err
		}
		res.Liked = liked
	} else {
		res.Liked = false
	}
	if stats.GetChannelStreaming(r.ID) {
		res.HLS = fmt.Sprintf("/app/channel%sstream/llhls.m3u8", r.ID)
	} else {
		res.HLS = fmt.Sprintf("/app/channel%s/llhls.m3u8", r.ID)
	}

	return res, nil
}

func (s *ChannelServiceImpl) GetSchedule(r requests.GetScheduleRequest) (responses.GetScheduleResponse, error) {
	var res responses.GetScheduleResponse
	pastTracks, err := repository.NewScheduleDB().GetPastTracks(r.ChannelID, r.Past)
	if err != nil {
		return res, err
	}
	res.Past = pastTracks

	nextTracks, err := repository.NewScheduleDB().GetNextTracks(r.ChannelID, r.Next)
	if err != nil {
		return res, err
	}
	res.Next = nextTracks

	return res, nil
}

func (s *ChannelServiceImpl) GetScheduleRange(r requests.GetScheduleRangeRequest) (responses.GetScheduleRangeResponse, error) {
	var res responses.GetScheduleRangeResponse
	tracks, err := repository.NewScheduleDB().GetTracksInRange(r.ChannelID, r.From, r.To)
	if err != nil {
		return res, err
	}
	res.Tracks = tracks

	return res, nil
}

func (s *ChannelServiceImpl) StartStream(r requests.StartStreamRequest) error {
	stats.SetChannelStreaming(r.ID, true)
	return nil
}

func (s *ChannelServiceImpl) StopStream(r requests.StopStreamRequest) error {
	stats.SetChannelStreaming(r.ID, false)
	return nil
}
