package repository

import (
	"database/sql"
	"fmt"
	"netradio/internal/model"
	"netradio/pkg/database"
	"netradio/pkg/files"
	"strings"
	"time"
)

type ScheduleDB interface {
	AddTracksToSchedule([]model.ScheduleTrack) error
	GetPastTracks(id string, count int) ([]model.ScheduleTrackFull, error)
	GetNextTracks(id string, count int) ([]model.ScheduleTrackFull, error)
	GetTracksInRange(id string, from, to time.Time) ([]model.ScheduleTrackFull, error)
	DeleteTrack(id string) error
	UpdateTracks(tracks []model.ScheduleTrack) error
}

func NewScheduleDB() ScheduleDB {
	return &ScheduleDBImpl{
		conn: database.GetConnection(),
	}
}

type ScheduleDBImpl struct {
	conn *sql.DB
}

func (db *ScheduleDBImpl) AddTracksToSchedule(tracks []model.ScheduleTrack) error {
	tx, err := db.conn.Begin()

	valueStrings := []string{}
	valueArgs := []interface{}{}
	for i, track := range tracks {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d::integer, $%d::integer, $%d::timestamp, $%d::timestamp)", i*4+1, i*4+2, i*4+3, i*4+4))
		valueArgs = append(valueArgs, track.ChannelId)
		valueArgs = append(valueArgs, track.TrackId)
		valueArgs = append(valueArgs, track.StartDate)
		valueArgs = append(valueArgs, track.EndDate)
	}

	stmt := fmt.Sprintf("INSERT INTO schedule (channelid, trackid, startdate, enddate) VALUES %s", strings.Join(valueStrings, ","))
	_, err = tx.Exec(stmt, valueArgs...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()

	return err
}

func (db *ScheduleDBImpl) DeleteTrack(id string) error {
	_, err := db.conn.Exec("DELETE FROM schedule WHERE id=$1", id)
	return err
}

func (db *ScheduleDBImpl) UpdateTracks(tracks []model.ScheduleTrack) error {
	tx, err := db.conn.Begin()

	valueStrings := []string{}
	valueArgs := []interface{}{}
	for i, track := range tracks {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d::integer, $%d::integer, $%d::timestamp, $%d::timestamp)", i*4+1, i*4+2, i*4+3, i*4+4))
		valueArgs = append(valueArgs, track.ID)
		valueArgs = append(valueArgs, track.TrackId)
		valueArgs = append(valueArgs, track.StartDate)
		valueArgs = append(valueArgs, track.EndDate)
	}

	stmt := fmt.Sprintf("UPDATE schedule AS s1 SET trackid=s2.trackid, startdate=s2.startdate, enddate=s2.enddate FROM (VALUES %s) as s2(id, trackid, startdate, enddate) WHERE s1.id=s2.id", strings.Join(valueStrings, ","))
	_, err = tx.Exec(stmt, valueArgs...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()

	return err
}

func (db *ScheduleDBImpl) GetPastTracks(id string, count int) ([]model.ScheduleTrackFull, error) {
	rows, err := db.conn.Query("SELECT schedule.id, schedule.channelid, schedule.startdate, schedule.enddate, tracks.id, tracks.title, tracks.performancer, tracks.year, tracks.audio, tracks.duration FROM tracks JOIN schedule ON schedule.trackid=tracks.id WHERE schedule.channelid=$1 AND schedule.startdate < NOW() ORDER BY schedule.startdate DESC LIMIT $2", id, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return ScanScheduleRows(rows)
}

func (db *ScheduleDBImpl) GetNextTracks(id string, count int) ([]model.ScheduleTrackFull, error) {
	rows, err := db.conn.Query("SELECT schedule.id, schedule.channelid, schedule.startdate, schedule.enddate, tracks.id, tracks.title, tracks.performancer, tracks.year, tracks.audio, tracks.duration FROM tracks JOIN schedule ON schedule.trackid=tracks.id WHERE schedule.channelid=$1 AND schedule.startdate > NOW() ORDER BY schedule.startdate ASC LIMIT $2", id, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return ScanScheduleRows(rows)
}

func (db *ScheduleDBImpl) GetTracksInRange(id string, from, to time.Time) ([]model.ScheduleTrackFull, error) {
	rows, err := db.conn.Query("SELECT schedule.id, schedule.channelid, schedule.startdate, schedule.enddate, tracks.id, tracks.title, tracks.performancer, tracks.year, tracks.audio, tracks.duration FROM tracks JOIN schedule ON schedule.trackid=tracks.id WHERE schedule.channelid=$1 AND schedule.startdate > $2 AND schedule.enddate < $3 ORDER BY schedule.startdate ASC", id, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return ScanScheduleRows(rows)
}

func ScanScheduleRows(rows *sql.Rows) ([]model.ScheduleTrackFull, error) {
	res := make([]model.ScheduleTrackFull, 0)
	for rows.Next() {
		var temp model.ScheduleTrackFull
		var track model.Track
		var audio sql.NullString
		var duration sql.NullInt64
		err := rows.Scan(&temp.ID, &temp.ChannelId, &temp.StartDate, &temp.EndDate, &track.ID, &track.Title, &track.Performancer, &track.Year, &audio, &duration)
		if err != nil {
			return res, err
		}
		if audio.Valid {
			track.Audio = files.ToURL(audio.String)
		}
		if duration.Valid {
			track.Duration = time.Duration(duration.Int64)
		}
		temp.Track = track
		res = append(res, temp)
	}

	return res, nil
}
