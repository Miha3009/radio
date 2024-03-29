package service

import (
	"netradio/internal/controller/requests"
	"netradio/internal/controller/responses"
	"netradio/internal/model"
	"netradio/internal/repository"
	"netradio/pkg/stats"
	"strconv"
	"time"
)

type TrackService interface {
	GetTrack(r requests.GetTrackRequest) (responses.GetTrackResponse, error)
	GetTrackList(r requests.GetTrackListRequest) (responses.GetTrackListResponse, error)
	CreateTrack(r requests.CreateTrackRequest) (responses.CreateTrackResponse, error)
	DeleteTrack(r requests.DeleteTrackRequest) error
	UpdateTrack(r requests.UpdateTrackRequest) (responses.UpdateTrackResponse, error)
	LikeTrack(r requests.LikeTrackRequest) error
	GetTrackComments(r requests.GetTrackCommentsRequest) (responses.GetTrackCommentsResponse, error)
	CommentTrack(r requests.CommentTrackRequest) (responses.CommentTrackResponse, error)
	UploadTrack(r requests.UploadTrackRequest) error
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
	res.ID = track.ID
	res.Title = track.Title
	res.Performancer = track.Performancer
	res.Year = track.Year
	res.Duration = track.Duration
	res.Audio = track.Audio

	likeCount, err := s.db.LikeCount(r.ID)
	if err != nil {
		return res, err
	}
	res.LikeCount = likeCount

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

func (s *TrackServiceImpl) GetTrackList(r requests.GetTrackListRequest) (responses.GetTrackListResponse, error) {
	var res responses.GetTrackListResponse
	tracks, tracksCount, err := s.db.GetTrackList(r.Offset, r.Limit, r.Query)
	if err != nil {
		return res, err
	}
	res.Tracks = tracks
	res.Count = tracksCount

	return res, nil
}

func (s *TrackServiceImpl) CreateTrack(r requests.CreateTrackRequest) (responses.CreateTrackResponse, error) {
	var res responses.CreateTrackResponse
	var track model.Track
	track.Title = r.Title
	track.Performancer = r.Performancer
	track.Year = r.Year

	id, err := s.db.CreateTrack(track)
	if err != nil {
		return res, err
	}
	res.ID = strconv.Itoa(id)

	return res, nil
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
	if r.Performancer != nil {
		track.Performancer = *r.Performancer
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
		stats.AddLike(r.ID)
		return s.db.LikeTrack(r.ID, strconv.Itoa(r.UserID))
	} else {
		stats.RemoveLike(r.ID)
		return s.db.UnlikeTrack(r.ID, strconv.Itoa(r.UserID))
	}
}

func (s *TrackServiceImpl) GetTrackComments(r requests.GetTrackCommentsRequest) (responses.GetTrackCommentsResponse, error) {
	var res responses.GetTrackCommentsResponse
	comments, err := s.db.GetTrackComments(r.ID)
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

func (s *TrackServiceImpl) CommentTrack(r requests.CommentTrackRequest) (responses.CommentTrackResponse, error) {
	var res responses.CommentTrackResponse
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

	return res, s.db.CommentTrack(r.ID, strconv.Itoa(commentId))
}

func (s *TrackServiceImpl) UploadTrack(r requests.UploadTrackRequest) error {
	return s.db.ChangeTrackAudio(r.ID, r.Audio, 0)
}
