package adminpanel

import (
	"io"
	"netradio/internal/repository"
	"netradio/pkg/context"
	"netradio/pkg/handlers"
	"os"
)

func newChangePodcastHandler(musicServise repository.MusicDB) *changePodcastHandler {
	return &changePodcastHandler{
		musicServise: musicServise,
	}
}

type changePodcastHandler struct {
	musicServise repository.MusicDB
}

func (h *changePodcastHandler) ServeHTTP(context context.Context) (handlers.Response, error) {
	request := context.GetRequest()
	err := request.ParseMultipartForm(maxPodcastSize)
	if err != nil {
		return handlers.Response{}, err
	}

	file, _, err := request.FormFile("music_file")
	if err != nil {
		return handlers.Response{}, err
	}
	content, err := io.ReadAll(file)
	if err != nil {
		return handlers.Response{}, err
	}
	err = os.WriteFile("new_music", content, os.ModePerm)

	request.PostFormValue("id")
	request.PostFormValue("title")
	request.PostFormValue("photo_link") // ?

	if err != nil {
		return handlers.Response{}, err
	}

	return handlers.Response{}, nil
}
