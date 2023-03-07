package adminpanel

import (
	"encoding/json"
	"errors"
	"netradio/internal/model"
	"netradio/internal/repository"
	"netradio/pkg/context"
	"netradio/pkg/handlers"
	"time"
)

func newChangeNewsHandler(newsService repository.NewsDB) *changeNewsHandler {
	return &changeNewsHandler{
		newsService: newsService,
	}
}

type changeNewsHandler struct {
	newsService repository.NewsDB
}

func (h *changeNewsHandler) ServeHTTP(context context.Context) (handlers.Response, error) {
	var rawNews struct {
		ID      *int    `json:"id,omitempty"`
		Title   *string `json:"title,omitempty"`
		Content *string `json:"content,omitempty"`
		//PublicationTime *int64  `json:"publication_date,omitempty"` // <- not allowed to modify
	}

	decoder := json.NewDecoder(context.GetRequest().Body)
	err := decoder.Decode(&rawNews)

	if err != nil {
		return handlers.Response{}, err
	}
	if rawNews.ID == nil {
		return handlers.Response{}, errors.New("no such have been passed")
	}

	modifiedNews, err := h.newsService.Get(*rawNews.ID)
	if err != nil {
		modifiedNews = model.News{
			ID: *rawNews.ID,
		}
	}

	if rawNews.Title != nil {
		modifiedNews.Title = *rawNews.Title
	}

	if rawNews.Content != nil {
		modifiedNews.Content = *rawNews.Content
	}

	modifiedNews.PublicationTime = time.Now().Unix()
	h.newsService.Add(modifiedNews)

	return handlers.Response{}, nil
}
