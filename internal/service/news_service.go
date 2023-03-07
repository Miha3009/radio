package service

import (
	"netradio/internal/model"
	"netradio/internal/repository"
)

type NewsService interface {
	Get(id int) (model.News, error)
	GetRange(offset, limit int) ([]model.News, error)
}

func NewNewsService() NewsService {
	return &NewsServiceImpl{
		db: repository.NewNewsDB(),
	}
}

type NewsServiceImpl struct {
	db repository.NewsDB
}

func (s *NewsServiceImpl) Get(id int) (model.News, error) {
	return s.db.Get(id)
}

func (s *NewsServiceImpl) GetRange(offset, limit int) ([]model.News, error) {
	return s.db.GetRange(offset, limit)
}
