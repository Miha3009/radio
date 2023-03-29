package hls

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"netradio/internal/model"
	"netradio/internal/repository"
	"netradio/pkg/log"
	"netradio/pkg/stats"
	"os"
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

	f, err := os.Create(channelID + ".txt")
	if err != nil {
		log.NewLogger().Error(err)
		return
	}
	defer f.Close()

	for i := 0; i < 1000; i += 1 {
		_, err := f.WriteString(fmt.Sprintf("file '%s.%d'\n", channelID, i%2))
		if err != nil {
			log.NewLogger().Error(err)
			return
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var startTime time.Time
	var duration time.Duration
	for i := 0; i < 1000; {
		if stats.GetChannelStatus(channelID) != model.ActiveChannel {
			return
		}

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

		out, err := os.Create(fmt.Sprintf("%s.%d", channelID, i%2))
		if err != nil {
			log.NewLogger().Error(err)
			continue
		}

		resp, err := http.Get(track.Audio)
		if err != nil {
			log.NewLogger().Error(err)
			out.Close()
			continue
		}

		_, err = io.Copy(out, resp.Body)
		out.Close()
		resp.Body.Close()
		if err != nil {
			log.NewLogger().Error(err)
			continue
		}

		if i == 0 {
			go func(ctx context.Context) {
				cmd := exec.CommandContext(ctx, "ffmpeg", "-re", "-f", "concat", "-i", fmt.Sprintf("%s.txt", channelID), "-c:a", "aac", "-ar", "44100", "-ac", "1", "-f", "flv", fmt.Sprintf("rtmp://localhost:1935/app/channel%s", channelID))
				_, err = cmd.Output()
				if err != nil {
					log.NewLogger().Error(err)
				}
			}(ctx)
			startTime = time.Now()
		} else {
			time.Sleep(duration - time.Now().Sub(startTime) + time.Second*10)
			startTime = startTime.Add(duration)
		}
		i += 1
		duration = track.Duration
	}
}

func GetCurrentTrackTime(channelID string) time.Duration {
	return *channelsToTrackTime[channelID]
}
