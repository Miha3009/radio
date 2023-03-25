package stats

import (
	"netradio/internal/repository"
	"netradio/pkg/log"
	"sync"
)

var (
	listeners           map[string]int
	listenerLock        sync.RWMutex
	trackLikes          map[string]int
	trackLock           sync.RWMutex
	trackForChannel     map[string]string
	trackForChannelLock sync.RWMutex
)

func Init() {
	listeners = make(map[string]int, 0)
	listenerLock = sync.RWMutex{}
	trackLikes = make(map[string]int, 0)
	trackLock = sync.RWMutex{}
	trackForChannel = make(map[string]string, 0)
	trackForChannelLock = sync.RWMutex{}

	tracks, likes, err := repository.NewTrackDB().GetLikesList()
	if err != nil {
		log.NewLogger().Error(err)
		return
	}
	for i := range tracks {
		trackLikes[tracks[i]] = likes[i]
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
