package service

import (
	"netradio/internal/controller/requests"
	"netradio/internal/controller/responses"
	"netradio/internal/model"
	"netradio/internal/repository"
	"strconv"
	"time"
)

type NewsService interface {
	GetNewsList(r requests.GetNewsListRequest) (responses.GetNewsListResponse, error)
	GetNews(r requests.GetNewsRequest) (responses.GetNewsResponse, error)
	CreateNews(r requests.CreateNewsRequest) (responses.CreateNewsResponse, error)
	DeleteNews(r requests.DeleteNewsRequest) error
	UpdateNews(r requests.UpdateNewsRequest) (responses.UpdateNewsResponse, error)
}

func NewNewsService() NewsService {
	return &NewsServiceImpl{
		db: repository.NewNewsDB(),
	}
}

type NewsServiceImpl struct {
	db repository.NewsDB
}

func (s *NewsServiceImpl) GetNewsList(r requests.GetNewsListRequest) (responses.GetNewsListResponse, error) {
	var res responses.GetNewsListResponse
	newsList, err := s.db.GetNewsList(r.Offset, r.Limit, r.Query)
	if err != nil {
		return res, err
	}

	newsCount, err := s.db.GetNewsCount()
	if err != nil {
		return res, err
	}

	res.NewsList = newsList
	res.Count = newsCount

	return res, nil
}

func (s *NewsServiceImpl) GetNews(r requests.GetNewsRequest) (responses.GetNewsResponse, error) {
	var res responses.GetNewsResponse
	news, err := s.db.GetNewsById(r.ID)
	if err != nil {
		return res, err
	}
	if news == nil {
		res.Found = false
		return res, err
	}

	res.Found = true
	res.ID = strconv.Itoa(news.ID)
	res.Title = news.Title
	res.Content = news.Content
	res.PublicationDate = news.PublicationDate

	return res, nil
}

func (s *NewsServiceImpl) CreateNews(r requests.CreateNewsRequest) (responses.CreateNewsResponse, error) {
	var res responses.CreateNewsResponse
	var news model.News
	news.Title = r.Title
	news.Content = r.Content
	news.PublicationDate = time.Now()

	id, err := s.db.CreateNews(news)
	if err != nil {
		return res, err
	}
	res.ID = strconv.Itoa(id)

	return res, nil
}

func (s *NewsServiceImpl) DeleteNews(r requests.DeleteNewsRequest) error {
	return s.db.DeleteNews(r.ID)
}

func (s *NewsServiceImpl) UpdateNews(r requests.UpdateNewsRequest) (responses.UpdateNewsResponse, error) {
	var res responses.UpdateNewsResponse
	news, err := s.db.GetNewsById(r.ID)
	if err != nil {
		return res, err
	}
	if news == nil {
		res.Found = false
		return res, nil
	}

	res.Found = true
	if r.Title != nil {
		news.Title = *r.Title
	}
	if r.Content != nil {
		news.Content = *r.Content
	}

	err = s.db.UpdateNews(*news)
	if err != nil {
		return res, err
	}

	return res, nil
}
