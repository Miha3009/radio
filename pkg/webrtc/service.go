package webrtc

import (
	"errors"
	"io"
	"netradio/internal/repository"
	"netradio/pkg/log"
	"os"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/pion/webrtc/v3/pkg/media/oggreader"
)

const oggPageDuration = time.Millisecond * 20

var (
	channelsToTracks map[string]*webrtc.TrackLocalStaticSample
	config           webrtc.Configuration
)

func StartAllChannels() {
	config = webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	channelsToTracks = make(map[string]*webrtc.TrackLocalStaticSample)

	channels, err := repository.NewChannelDB().GetChannels()
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

	channelsToTracks[channelID] = audioTrack

	for {
		path, err := repository.NewChannelDB().GetCurrentTrack(channelID)
		if err != nil {
			log.NewLogger().Error(err)
			return
		}

		file, err := os.Open(path)
		if err != nil {
			log.NewLogger().Error(err)
			return
		}

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
