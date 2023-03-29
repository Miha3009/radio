package stats

import (
	"netradio/internal/model"
	"netradio/internal/repository"
	"netradio/pkg/log"
	"sync"
	"time"
)

var (
	listeners            map[string]int
	listenerLock         sync.RWMutex
	trackLikes           map[string]int
	trackLock            sync.RWMutex
	trackForChannel      map[string]string
	trackForChannelLock  sync.RWMutex
	channelStatuses      map[string]model.ChannelStatus
	channelStatusesLock  sync.RWMutex
	channelStreaming     map[string]bool
	channelStreamingLock sync.RWMutex
)

func Init() {
	listeners = make(map[string]int, 0)
	listenerLock = sync.RWMutex{}
	trackLikes = make(map[string]int, 0)
	trackLock = sync.RWMutex{}
	trackForChannel = make(map[string]string, 0)
	trackForChannelLock = sync.RWMutex{}
	channelStatuses = make(map[string]model.ChannelStatus, 0)
	channelStatusesLock = sync.RWMutex{}
	channelStreaming = make(map[string]bool, 0)
	channelStreamingLock = sync.RWMutex{}

	tracks, likes, err := repository.NewTrackDB().GetLikesList()
	if err != nil {
		log.NewLogger().Error(err)
		return
	}
	for i := range tracks {
		trackLikes[tracks[i]] = likes[i]
	}

	channels, _, err := repository.NewChannelDB().GetChannels(0, 1000000, "", "")
	if err != nil {
		log.NewLogger().Fatal(err)
	}

	for _, channel := range channels {
		SetChannelStatus(channel.ID, channel.Status)
		go RunForChannel(channel.ID)
	}
}

func AddListener(channelID string) {
	listenerLock.Lock()
	defer listenerLock.Unlock()
	if count, ok := listeners[channelID]; ok {
		listeners[channelID] = count + 1
	} else {
		listeners[channelID] = 1
	}
}

func RemoveListener(channelID string) {
	listenerLock.Lock()
	defer listenerLock.Unlock()
	if count, ok := listeners[channelID]; ok && count > 0 {
		listeners[channelID] = count - 1
	} else {
		listeners[channelID] = 0
	}
}

func GetListeners(channelID string) int {
	listenerLock.RLock()
	defer listenerLock.RUnlock()
	if count, ok := listeners[channelID]; ok {
		return count
	} else {
		return 0
	}
}

func AddLike(trackID string) {
	trackLock.Lock()
	defer trackLock.Unlock()
	if count, ok := trackLikes[trackID]; ok {
		trackLikes[trackID] = count + 1
	} else {
		trackLikes[trackID] = 0
	}
}

func RemoveLike(trackID string) {
	trackLock.Lock()
	defer trackLock.Unlock()
	if count, ok := trackLikes[trackID]; ok && count > 0 {
		trackLikes[trackID] = count - 1
	} else {
		trackLikes[trackID] = 0
	}
}

func GetLikes(channelID string) int {
	trackLock.RLock()
	defer trackLock.RUnlock()
	if count, ok := trackLikes[channelID]; ok {
		return count
	} else {
		return 0
	}
}

func SetTrackForChannel(trackID, channelID string) {
	trackForChannelLock.Lock()
	defer trackForChannelLock.Unlock()
	trackForChannel[channelID] = trackID
}

func GetTrackForChannel(channelID string) string {
	trackForChannelLock.RLock()
	defer trackForChannelLock.RUnlock()
	if res, ok := trackForChannel[channelID]; ok {
		return res
	} else {
		return ""
	}
}

func SetChannelStatus(channelID string, status model.ChannelStatus) {
	channelStatusesLock.Lock()
	defer channelStatusesLock.Unlock()
	channelStatuses[channelID] = status
}

func GetChannelStatus(channelID string) model.ChannelStatus {
	channelStatusesLock.RLock()
	defer channelStatusesLock.RUnlock()
	if status, ok := channelStatuses[channelID]; ok {
		return status
	} else {
		return model.StoppedChannel
	}
}

func RunForChannel(channelID string) {
	ticker := time.NewTicker(time.Second * 10)
	for ; true; <-ticker.C {
		if GetChannelStatus(channelID) == model.DeletedChannel {
			return
		}
		repository.NewStatsDB().AddListenerTimestamp(channelID, GetListeners(channelID))
	}
}

func SetChannelStreaming(channelID string, streaming bool) {
	channelStreamingLock.Lock()
	defer channelStreamingLock.Unlock()
	channelStreaming[channelID] = streaming
}

func GetChannelStreaming(channelID string) bool {
	channelStreamingLock.RLock()
	defer channelStreamingLock.RUnlock()
	if streaming, ok := channelStreaming[channelID]; ok {
		return streaming
	} else {
		return false
	}
}
