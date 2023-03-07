package adminpanel

import (
	"io"
	"netradio/internal/repository"
	"netradio/pkg/context"
	"netradio/pkg/handlers"
	"os"
)

const maxPodcastSize = 1 << 30 // 1 gigabyte

func newCreatePodcastHandler(musicService repository.MusicDB) *createHandler {
	return &createHandler{
		musicService: musicService,
	}
}

type createHandler struct {
	musicService repository.MusicDB
}

func (h *createHandler) ServeHTTP(context context.Context) (handlers.Response, error) {
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

	//request.PostFormValue("id") // no id when creation yet
	request.PostFormValue("title")
	request.PostFormValue("photo_link") // ?

	if err != nil {
		return handlers.Response{}, err
	}

	return handlers.Response{}, nil
}
