package service

import (
	"netradio/internal/controller/requests"
	"netradio/internal/controller/responses"
	"netradio/internal/model"
	"netradio/internal/repository"
	"strconv"
	"time"
)

type TrackService interface {
	GetTrack(r requests.GetTrackRequest) (responses.GetTrackResponse, error)
	CreateTrack(r requests.CreateTrackRequest) error
	DeleteTrack(r requests.DeleteTrackRequest) error
	UpdateTrack(r requests.UpdateTrackRequest) (responses.UpdateTrackResponse, error)
	LikeTrack(r requests.LikeTrackRequest) error
	CommentTrack(r requests.CommentTrackRequest) error
}

func NewTrackService() TrackService {
	return &TrackServiceImpl{
		db: repository.NewTrackDB(),
	}
}

type TrackServiceImpl struct {
	db repository.TrackDB
}

func (s *TrackServiceImpl) GetTrack(r requests.GetTrackRequest) (responses.GetTrackResponse, error) {
	var res responses.GetTrackResponse
	track, err := s.db.GetTrackById(r.ID)
	if err != nil {
		return res, err
	}
	if track == nil {
		res.Found = false
		return res, err
	}

	res.Found = true
	res.ID = strconv.Itoa(track.ID)
	res.Title = track.Title
	res.Perfomancer = track.Perfomancer
	res.Year = track.Year

	if r.UserID != nil {
		liked, err := s.db.IsTrackLiked(r.ID, *r.UserID)
		if err != nil {
			return res, err
		}
		res.Liked = liked
	} else {
		res.Liked = false
	}

	return res, nil
}

func (s *TrackServiceImpl) CreateTrack(r requests.CreateTrackRequest) error {
	var track model.Track
	track.Title = r.Title
	track.Perfomancer = r.Perfomancer
	track.Year = r.Year
	return s.db.CreateTrack(track)
}

func (s *TrackServiceImpl) DeleteTrack(r requests.DeleteTrackRequest) error {
	return s.db.DeleteTrack(r.ID)
}

func (s *TrackServiceImpl) UpdateTrack(r requests.UpdateTrackRequest) (responses.UpdateTrackResponse, error) {
	var res responses.UpdateTrackResponse
	track, err := s.db.GetTrackById(r.ID)
	if err != nil {
		return res, err
	}
	if track == nil {
		res.Found = false
		return res, nil
	}

	res.Found = true
	if r.Title != nil {
		track.Title = *r.Title
	}
	if r.Perfomancer != nil {
		track.Perfomancer = *r.Perfomancer
	}
	if r.Year != nil {
		track.Year = *r.Year
	}

	err = s.db.UpdateTrack(*track)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *TrackServiceImpl) LikeTrack(r requests.LikeTrackRequest) error {
	if r.Like {
		return s.db.LikeTrack(r.ID, strconv.Itoa(r.UserID))
	} else {
		return s.db.UnlikeTrack(r.ID, strconv.Itoa(r.UserID))
	}
}

func (s *TrackServiceImpl) CommentTrack(r requests.CommentTrackRequest) error {
	var comment model.Comment
	comment.UserID = r.UserID
	comment.Text = r.Text
	comment.Date = time.Now()
	if r.Parent != nil {
		parent, err := strconv.Atoi(*r.Parent)
		if err != nil {
			return err
		}
		comment.Parent = &parent
	}

	commentId, err := repository.NewCommentDB().CreateComment(comment)
	if err != nil {
		return err
	}

	return s.db.CommentTrack(r.ID, strconv.Itoa(commentId))
}
