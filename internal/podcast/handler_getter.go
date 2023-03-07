package podcast

import (
	"encoding/json"
	"netradio/internal/repository"
	"netradio/pkg/context"
)

func newGetterHandler(musicService repository.MusicDB) *getterHandler {
	return &getterHandler{
		musicService: musicService,
	}
}

type getterHandler struct {
	musicService repository.MusicDB
}

func (h *getterHandler) ServeJSON(context context.Context) (json.RawMessage, error) {
	podcasts := h.musicService.GetPodcasts()

	return json.Marshal(podcasts)
}
