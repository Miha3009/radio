package service

import (
	"errors"
	"io"
	"netradio/internal/controller/requests"
	"netradio/internal/controller/responses"
	"netradio/internal/model"
	"netradio/internal/repository"
	"os"
	"strconv"
	"time"

	"github.com/pion/webrtc/v3/pkg/media/oggreader"
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
	res.ID = strconv.Itoa(track.ID)
	res.Title = track.Title
	res.Perfomancer = track.Perfomancer
	res.Year = track.Year
	res.Duration = track.Duration

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
	tracks, err := s.db.GetTrackList(r.Offset, r.Limit, r.Query)
	if err != nil {
		return res, err
	}
	res.Tracks = tracks

	tracksCount, err := s.db.GetTracksCount()
	if err != nil {
		return res, err
	}
	res.Count = tracksCount

	return res, nil
}

func (s *TrackServiceImpl) CreateTrack(r requests.CreateTrackRequest) (responses.CreateTrackResponse, error) {
	var res responses.CreateTrackResponse
	var track model.Track
	track.Title = r.Title
	track.Perfomancer = r.Perfomancer
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
	file, err := os.Open(r.Audio)
	if err != nil {
		return err
	}
	defer file.Close()

	ogg, _, err := oggreader.NewWith(file)
	if err != nil {
		return err
	}

	var lastGranule uint64
	duration := time.Duration(0)

	for {
		_, pageHeader, err := ogg.ParseNextPage()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		sampleCount := float64(pageHeader.GranulePosition - lastGranule)
		lastGranule = pageHeader.GranulePosition
		duration = duration + time.Duration((sampleCount/48000)*1000)*time.Millisecond
	}

	return s.db.ChangeTrackAudio(r.ID, r.Audio, duration)
}
