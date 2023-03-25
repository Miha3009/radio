package webrtc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"netradio/internal/repository"
	"netradio/pkg/log"
	"netradio/pkg/stats"
	"os"
	"strconv"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/pion/webrtc/v3/pkg/media/oggreader"
	"golang.org/x/net/websocket"
)

const oggPageDuration = time.Millisecond * 20

var (
	channelsToTracks    map[string]*webrtc.TrackLocalStaticSample
	channelsToTrackTime map[string]*time.Duration
	config              webrtc.Configuration
)

func StartAllChannels() {
	config = webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs:           []string{"turn:relay.metered.ca:443"},
				Username:       "72572597ae7d4797af60789a",
				Credential:     "sWrcoUXDRy2SqIly",
				CredentialType: webrtc.ICECredentialTypePassword,
			},
		},
	}

	channelsToTracks = make(map[string]*webrtc.TrackLocalStaticSample)
	channelsToTrackTime = make(map[string]*time.Duration)

	channels, err := repository.NewChannelDB().GetChannels(0, 1000000, "", "")
	if err != nil {
		log.NewLogger().Fatal(err)
	}

	for _, channel := range channels {
		go StartChannel(channel.ID)
	}
}

func StartChannel(channelID string) {
	audioTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus}, "audio", "pion")
	if err != nil {
		log.NewLogger().Warn(err)
		return
	}

	currentTime := time.Second
	channelsToTracks[channelID] = audioTrack
	channelsToTrackTime[channelID] = &currentTime

	for {
		track, err := repository.NewChannelDB().GetCurrentTrack(channelID)
		if err != nil {
			log.NewLogger().Error(err)
			time.Sleep(time.Second)
			continue
		}
		if track == nil {
			time.Sleep(time.Second)
			continue
		}

		stats.SetTrackForChannel(strconv.Itoa(track.ID), channelID)
		file, err := os.Open(track.Audio)
		if err != nil {
			log.NewLogger().Error(err)
			time.Sleep(time.Second)
			continue
		}
		defer file.Close()

		ogg, _, err := oggreader.NewWith(file)
		if err != nil {
			log.NewLogger().Error(err)
		}

		var lastGranule uint64

		ticker := time.NewTicker(oggPageDuration)
		for ; true; <-ticker.C {
			pageData, pageHeader, err := ogg.ParseNextPage()
			if errors.Is(err, io.EOF) {
				break
			}

			if err != nil {
				log.NewLogger().Error(err)
			}

			sampleCount := float64(pageHeader.GranulePosition - lastGranule)
			lastGranule = pageHeader.GranulePosition
			sampleDuration := time.Duration((sampleCount/48000)*1000) * time.Millisecond
			currentTime = time.Duration((float64(pageHeader.GranulePosition)/48000)*1000) * time.Millisecond

			if err = audioTrack.WriteSample(media.Sample{Data: pageData, Duration: sampleDuration}); err != nil {
				log.NewLogger().Error(err)
			}
		}
		file.Close()
	}
}

func GetPeerConfig() webrtc.Configuration {
	return config
}

func GetAudioTrack(channelID string) (*webrtc.TrackLocalStaticSample, error) {
	if track, ok := channelsToTracks[channelID]; ok {
		return track, nil
	} else {
		return nil, errors.New("Track not found")
	}
}

func GetCurrentTrackTime(channelID string) time.Duration {
	return *channelsToTrackTime[channelID]
}

type ConnectStruct struct {
	ID string
}

type DataStruct struct {
	Listeners int `json:"listeners"`
	Likes     int `json:"likes"`
}

func (s *ConnectStruct) HandleConnectChannel(ws *websocket.Conn) {
	logger := log.NewLogger()
	peerConnection, err := webrtc.NewPeerConnection(GetPeerConfig())
	if err != nil {
		logger.Error(err)
		return
	}

	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}

		outbound, marshalErr := json.Marshal(c.ToJSON())
		if marshalErr != nil {
			logger.Error(marshalErr)
			return
		}

		if _, err = ws.Write(outbound); err != nil {
			logger.Error(err)
			return
		}
	})

	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		d.OnOpen(func() {
			stats.AddListener(s.ID)
			for range time.Tick(time.Second * 10) {
				res, err := json.Marshal(DataStruct{Listeners: stats.GetListeners(s.ID), Likes: stats.GetLikes(stats.GetTrackForChannel(s.ID))})
				if err != nil {
					logger.Error(err)
					return
				}
				if err = d.Send(res); err != nil {
					logger.Error(err)
					return
				}
			}
		})
		d.OnClose(func() {
			stats.RemoveListener(s.ID)
		})
	})

	track, err := GetAudioTrack(s.ID)
	if err != nil {
		logger.Error(err)
		return
	}

	rtpSender, err := peerConnection.AddTrack(track)
	if err != nil {
		logger.Error(err)
		return
	}

	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()

	for {
		buf := make([]byte, 10000)

		n, err := ws.Read(buf)
		if err != nil {
			logger.Error(err)
			return
		}

		var (
			candidate webrtc.ICECandidateInit
			offer     webrtc.SessionDescription
		)

		switch {
		case json.Unmarshal(buf[:n], &offer) == nil && offer.SDP != "":
			if err = peerConnection.SetRemoteDescription(offer); err != nil {
				logger.Error(err)
				return
			}

			answer, answerErr := peerConnection.CreateAnswer(nil)
			if answerErr != nil {
				logger.Error(answerErr)
				return
			}

			if err = peerConnection.SetLocalDescription(answer); err != nil {
				logger.Error(err)
				return
			}

			outbound, marshalErr := json.Marshal(answer)
			if marshalErr != nil {
				logger.Error(marshalErr)
				return
			}

			if _, err = ws.Write(outbound); err != nil {
				logger.Error(err)
				return
			}
		case json.Unmarshal(buf[:n], &candidate) == nil && candidate.Candidate != "":
			if err = peerConnection.AddICECandidate(candidate); err != nil {
				logger.Error(err)
				return
			}
		default:
			fmt.Println("Unknown message")
			return
		}
	}
}
