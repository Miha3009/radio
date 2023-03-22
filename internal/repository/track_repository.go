package repository

import (
	"database/sql"
	"netradio/internal/model"
	"netradio/pkg/database"
	"time"
)

type TrackDB interface {
	GetTracksCount() (int, error)
	GetTrackById(id string) (*model.Track, error)
	GetTrackList(offset, limit int, query string) ([]model.Track, error)
	CreateTrack(track model.Track) error
	UpdateTrack(track model.Track) error
	DeleteTrack(id string) error
	IsTrackLiked(id, userId string) (bool, error)
	LikeTrack(id, userId string) error
	UnlikeTrack(id, userId string) error
	CommentTrack(id, commentId string) error
	ChangeTrackAudio(id, audio string, duration time.Duration) error
	LikeCount(id string) (int, error)
}

func NewTrackDB() TrackDB {
	return &TrackDBImpl{
		conn: database.GetConnection(),
	}
}

type TrackDBImpl struct {
	conn *sql.DB
}

func (db *TrackDBImpl) GetTracksCount() (int, error) {
	count := 0
	err := db.conn.QueryRow("SELECT COUNT(*) FROM tracks").Scan(&count)
	return count, err
}

func (db *TrackDBImpl) GetTrackById(id string) (*model.Track, error) {
	var res model.Track
	rows, err := db.conn.Query("SELECT id, title, perfomancer, year, audio, duration FROM tracks WHERE id=$1", id)
	if err != nil {
		return &res, err
	}
	if rows.Next() {
		err = rows.Scan(&res.ID, &res.Title, &res.Perfomancer, &res.Year, &res.Audio, &res.Duration)
		return &res, err
	}

	return nil, nil
}

func (db *TrackDBImpl) GetTrackList(offset, limit int, query string) ([]model.Track, error) {
	res := make([]model.Track, 0)
	query = "%" + query + "%"
	rows, err := db.conn.Query("SELECT id, title, perfomancer, year, audio, duration FROM tracks WHERE title LIKE $3 OFFSET $1 LIMIT $2", offset, limit, query)
	if err != nil {
		return res, err
	}
	for rows.Next() {
		var temp model.Track
		var audio sql.NullString
		var duration sql.NullInt64
		err = rows.Scan(&temp.ID, &temp.Title, &temp.Perfomancer, &temp.Year, &audio, &duration)
		if err != nil {
			return res, err
		}
		if audio.Valid {
			temp.Audio = audio.String
		}
		if duration.Valid {
			temp.Duration = time.Duration(duration.Int64)
		}
		res = append(res, temp)
	}

	return res, nil
}

func (db *TrackDBImpl) CreateTrack(track model.Track) error {
	_, err := db.conn.Exec("INSERT INTO tracks (title, perfomancer, year) VALUES ($1, $2, $3)", track.Title, track.Perfomancer, track.Year)
	return err
}

func (db *TrackDBImpl) UpdateTrack(track model.Track) error {
	_, err := db.conn.Exec("UPDATE tracks SET title=$1, perfomancer=$2, year=$3 WHERE id=$4", track.Title, track.Perfomancer, track.Year, track.ID)
	return err
}

func (db *TrackDBImpl) DeleteTrack(id string) error {
	_, err := db.conn.Exec("DELETE FROM tracks WHERE id=$1", id)
	return err
}

func (db *TrackDBImpl) IsTrackLiked(id, userId string) (bool, error) {
	rows, err := db.conn.Query("SELECT * FROM tracks_likes WHERE trackid=$1 AND userid=$2", id, userId)
	if err != nil {
		return false, err
	}
	return rows.Next(), nil
}

func (db *TrackDBImpl) LikeTrack(id, userId string) error {
	_, err := db.conn.Exec("INSERT INTO tracks_likes (trackid, userid, time) VALUES ($1, $2, $3)", id, userId, time.Now())
	return err
}

func (db *TrackDBImpl) UnlikeTrack(id, userId string) error {
	_, err := db.conn.Exec("DELETE FROM tracks_likes WHERE trackid=$1 AND userid=$2", id, userId)
	return err
}

func (db *TrackDBImpl) CommentTrack(id, commentId string) error {
	_, err := db.conn.Exec("INSERT INTO tracks_comments (trackid, commentid) VALUES ($1, $2)", id, commentId)
	return err
}

func (db *TrackDBImpl) ChangeTrackAudio(id, audio string, duration time.Duration) error {
	_, err := db.conn.Exec("UPDATE tracks SET audio=$1, duration=$2 WHERE id=$3", audio, duration, id)
	return err
}

func (db *TrackDBImpl) LikeCount(id string) (int, error) {
	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM tracks_likes WHERE trackid=$1", id).Scan(&count)
	return count, err
}
