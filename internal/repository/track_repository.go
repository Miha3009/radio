package repository

import (
	"database/sql"
	"netradio/internal/model"
	"netradio/pkg/database"
	"netradio/pkg/files"
	"strconv"
	"time"
)

type TrackDB interface {
	GetTrackById(id string) (*model.Track, error)
	GetTrackList(offset, limit int, query string) ([]model.Track, int, error)
	CreateTrack(track model.Track) (int, error)
	UpdateTrack(track model.Track) error
	DeleteTrack(id string) error
	IsTrackLiked(id, userId string) (bool, error)
	LikeTrack(id, userId string) error
	UnlikeTrack(id, userId string) error
	GetTrackComments(id string) ([]model.Comment, error)
	CommentTrack(id, commentId string) error
	ChangeTrackAudio(id, audio string, duration time.Duration) error
	LikeCount(id string) (int, error)
	GetLikesList() ([]string, []int, error)
}

func NewTrackDB() TrackDB {
	return &TrackDBImpl{
		conn: database.GetConnection(),
	}
}

type TrackDBImpl struct {
	conn *sql.DB
}

func (db *TrackDBImpl) GetTrackById(id string) (*model.Track, error) {
	var res model.Track
	rows, err := db.conn.Query("SELECT id, title, performancer, year, audio, duration FROM tracks WHERE id=$1", id)
	if err != nil {
		return &res, err
	}
	defer rows.Close()
	if rows.Next() {
		var audio sql.NullString
		var duration sql.NullInt64
		err = rows.Scan(&res.ID, &res.Title, &res.Performancer, &res.Year, &audio, &duration)
		if audio.Valid {
			res.Audio = files.ToURL(audio.String)
		}
		if duration.Valid {
			res.Duration = time.Duration(duration.Int64)
		}
		return &res, err
	}

	return nil, nil
}

func (db *TrackDBImpl) GetTrackList(offset, limit int, query string) ([]model.Track, int, error) {
	res := make([]model.Track, 0)
	query = "%" + query + "%"
	rows, err := db.conn.Query("SELECT id, title, performancer, year, audio, duration FROM tracks WHERE title LIKE $3 OFFSET $1 LIMIT $2", offset, limit, query)
	if err != nil {
		return res, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var temp model.Track
		var audio sql.NullString
		var duration sql.NullInt64
		err = rows.Scan(&temp.ID, &temp.Title, &temp.Performancer, &temp.Year, &audio, &duration)
		if err != nil {
			return res, 0, err
		}
		if audio.Valid {
			temp.Audio = files.ToURL(audio.String)
		}
		if duration.Valid {
			temp.Duration = time.Duration(duration.Int64)
		}
		res = append(res, temp)
	}

	count := 0
	err = db.conn.QueryRow("SELECT COUNT(*) FROM tracks WHERE title LIKE $1", query).Scan(&count)

	return res, count, err
}

func (db *TrackDBImpl) CreateTrack(track model.Track) (int, error) {
	var id int
	err := db.conn.QueryRow("INSERT INTO tracks (title, performancer, year) VALUES ($1, $2, $3) RETURNING id", track.Title, track.Performancer, track.Year).Scan(&id)
	return id, err
}

func (db *TrackDBImpl) UpdateTrack(track model.Track) error {
	_, err := db.conn.Exec("UPDATE tracks SET title=$1, performancer=$2, year=$3 WHERE id=$4", track.Title, track.Performancer, track.Year, track.ID)
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
	defer rows.Close()
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

func (db *TrackDBImpl) GetTrackComments(id string) ([]model.Comment, error) {
	res := make([]model.Comment, 0)
	rows, err := db.conn.Query("SELECT comments.id, users.id, users.name, users.avatar, comments.parent, comments.text, comments.time FROM tracks_comments JOIN comments ON tracks_comments.commentid=comments.id JOIN users ON comments.userid=users.id WHERE tracks_comments.trackid=$1", id)
	if err != nil {
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		var temp model.Comment
		var parent sql.NullInt32
		var avatar sql.NullString
		err = rows.Scan(&temp.ID, &temp.UserID, &temp.UserName, &avatar, &parent, &temp.Text, &temp.Date)
		if err != nil {
			return res, err
		}
		if avatar.Valid {
			temp.UserAvatar = files.ToURL(avatar.String)
		}
		if parent.Valid {
			temp.Parent = int(parent.Int32)
		} else {
			temp.Parent = -1
		}
		res = append(res, temp)
	}

	return res, nil
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

func (db *TrackDBImpl) GetLikesList() ([]string, []int, error) {
	rows, err := db.conn.Query("SELECT trackid, COUNT(*) FROM tracks_likes GROUP BY trackid")
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	tracks := make([]string, 0)
	likes := make([]int, 0)
	for rows.Next() {
		like := 0
		track := 0
		err = rows.Scan(&track, &like)
		if err != nil {
			return tracks, likes, err
		}
		tracks = append(tracks, strconv.Itoa(track))
		likes = append(likes, like)
	}

	return tracks, likes, nil

}
