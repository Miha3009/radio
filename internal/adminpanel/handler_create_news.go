package adminpanel

import (
	"encoding/json"
	"netradio/internal/repository"
	"netradio/internal/model"
	"netradio/pkg/context"
	"netradio/pkg/handlers"
	"time"
)

func newCreateNewsHandler(newsService repository.NewsDB) *createNewsHandler {
	return &createNewsHandler{
		newsService: newsService,
	}
}

type createNewsHandler struct {
	newsService repository.NewsDB
}

func (h *createNewsHandler) ServeHTTP(context context.Context) (handlers.Response, error) {
	request := context.GetRequest()
	var rawNews struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&rawNews)
	if err != nil {
		return handlers.Response{}, err
	}

	newNews := model.News{
		Title:           rawNews.Title,
		Content:         rawNews.Content,
		PublicationTime: time.Now().Unix(),
	}

	h.newsService.Add(newNews)

	return handlers.Response{}, nil
}
