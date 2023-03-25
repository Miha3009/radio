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
	UploadImage(r requests.UploadImageRequest) (responses.UploadImageResponse, error)
	LikeNews(r requests.LikeNewsRequest) error
	GetNewsComments(r requests.GetNewsCommentsRequest) (responses.GetNewsCommentsResponse, error)
	CommentNews(r requests.CommentNewsRequest) (responses.CommentNewsResponse, error)
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
	res.Image = news.Image

	likeCount, err := s.db.LikeCount(r.ID)
	if err != nil {
		return res, err
	}
	res.LikeCount = likeCount

	if r.UserID != nil {
		liked, err := s.db.IsNewsLiked(r.ID, *r.UserID)
		if err != nil {
			return res, err
		}
		res.Liked = liked
	} else {
		res.Liked = false
	}

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

func (s *NewsServiceImpl) UploadImage(r requests.UploadImageRequest) (responses.UploadImageResponse, error) {
	var res responses.UploadImageResponse
	err := s.db.ChangeImage(r.ID, r.Image)
	if err != nil {
		return res, err
	}

	res.Image = r.Image

	return res, nil
}

func (s *NewsServiceImpl) LikeNews(r requests.LikeNewsRequest) error {
	if r.Like {
		return s.db.LikeNews(r.ID, strconv.Itoa(r.UserID))
	} else {
		return s.db.UnlikeNews(r.ID, strconv.Itoa(r.UserID))
	}
}

func (s *NewsServiceImpl) GetNewsComments(r requests.GetNewsCommentsRequest) (responses.GetNewsCommentsResponse, error) {
	var res responses.GetNewsCommentsResponse
	comments, err := s.db.GetNewsComments(r.ID)
	if err != nil {
		return res, err
	}
	commentMap := make(map[int]*model.Comment)
	for i := range comments {
		comments[i].Children = make([]model.Comment, 0)
		commentMap[comments[i].ID] = &comments[i]
	}
	for i := range comments {
		if comments[i].Parent != -1 {
			commentMap[comments[i].Parent].Children = append(commentMap[comments[i].Parent].Children, comments[i])
		}
	}
	newComments := make([]model.Comment, 0)
	for i := range comments {
		if comments[i].Parent == -1 {
			newComments = append(newComments, comments[i])
		}
	}
	res.Comments = newComments

	return res, nil
}

func (s *NewsServiceImpl) CommentNews(r requests.CommentNewsRequest) (responses.CommentNewsResponse, error) {
	var res responses.CommentNewsResponse
	var comment model.Comment
	comment.UserID = r.UserID
	comment.Text = r.Text
	comment.Date = time.Now()
	if r.Parent != nil {
		parent, err := strconv.Atoi(*r.Parent)
		if err != nil {
			return res, err
		}
		comment.Parent = parent
	}

	commentId, err := repository.NewCommentDB().CreateComment(comment)
	if err != nil {
		return res, err
	}
	res.ID = commentId

	return res, s.db.CommentNews(r.ID, strconv.Itoa(commentId))
}
