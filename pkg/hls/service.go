package hls

import (
	"fmt"
	"netradio/internal/repository"
	"netradio/pkg/log"
	"os/exec"
	"time"
)

var (
	channelsToTrackTime map[string]*time.Duration
	streaming           bool
)

func StartAllChannels(isStreaming bool) {
	channelsToTrackTime = make(map[string]*time.Duration)
	streaming = isStreaming

	channels, _, err := repository.NewChannelDB().GetChannels(0, 1000000, "", "")
	if err != nil {
		log.NewLogger().Fatal(err)
	}

	for _, channel := range channels {
		go StartChannel(channel.ID)
	}
}

func StartChannel(channelID string) {
	if !streaming {
		return
	}
	currentTime := time.Duration(0)
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

		cmd := exec.Command("ffmpeg", "-re", "-i", track.Audio, "-c:a", "aac", "-ar", "44100", "-ac", "1", "-f", "flv", fmt.Sprintf("rtmp://localhost:1935/app/channel%s", channelID))
		_, err = cmd.Output()
		if err != nil {
			log.NewLogger().Error(err)
			time.Sleep(time.Second)
			continue
		}
	}
}

func GetCurrentTrackTime(channelID string) time.Duration {
	return *channelsToTrackTime[channelID]
}
