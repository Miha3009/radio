package service

import (
	"netradio/internal/controller/requests"
	"netradio/internal/controller/responses"
	"netradio/internal/model"
	"netradio/internal/repository"
	webrtchelper "netradio/pkg/webrtc"
	"strconv"

	"github.com/pion/webrtc/v3"
)

type ChannelService interface {
	GetChannels() (responses.GetChannelsResponse, error)
	GetChannel(r requests.GetChannelRequest) (responses.GetChannelResponse, error)
	CreateChannel(r requests.CreateChannelRequest) error
	DeleteChannel(r requests.DeleteChannelRequest) error
	UpdateChannel(r requests.UpdateChannelRequest) (responses.UpdateChannelResponse, error)
	StartChannel(r requests.StartChannelRequest) error
	StopChannel(r requests.StopChannelRequest) error
	ConnectChannel(r requests.ConnectChannelRequest) (webrtc.SessionDescription, error)
	UploadLogo(r requests.UploadLogoRequest) error
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
	res.Status = int(channel.Status)

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

func (s *ChannelServiceImpl) StartChannel(r requests.StartChannelRequest) error {
	return s.db.ChangeChannelStatus(r.ID, model.ActiveChannel)
}

func (s *ChannelServiceImpl) StopChannel(r requests.StopChannelRequest) error {
	return s.db.ChangeChannelStatus(r.ID, model.StoppedChannel)
}

func (s *ChannelServiceImpl) ConnectChannel(r requests.ConnectChannelRequest) (webrtc.SessionDescription, error) {
	var offer webrtc.SessionDescription
	offer.SDP = r.SDP
	offer.Type = webrtc.SDPTypeOffer

	pc, err := webrtc.NewPeerConnection(webrtchelper.GetPeerConfig())
	if err != nil {
		return webrtc.SessionDescription{}, err
	}

	track, err := webrtchelper.GetAudioTrack(r.ID)
	if err != nil {
		return webrtc.SessionDescription{}, err
	}

	rtpSender, err := pc.AddTrack(track)
	if err != nil {
		return webrtc.SessionDescription{}, err
	}

	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()

	err = pc.SetRemoteDescription(offer)
	if err != nil {
		return webrtc.SessionDescription{}, err
	}

	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		return webrtc.SessionDescription{}, err
	}

	err = pc.SetLocalDescription(answer)
	if err != nil {
		return webrtc.SessionDescription{}, err
	}

	return answer, nil
}

func (s *ChannelServiceImpl) UploadLogo(r requests.UploadLogoRequest) error {
	return s.db.ChangeLogo(r.ID, r.Logo)
}
