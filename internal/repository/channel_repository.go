package repository

import (
	"database/sql"
	"errors"
	"netradio/internal/model"
	"netradio/pkg/database"
	"time"
)

type ChannelDB interface {
	GetChannels() ([]model.ChannelShortInfo, error)
	GetChannelById(id string) (*model.ChannelInfo, error)
	GetCurrentTrack(id string) (string, error)
	CreateChannel(channel model.ChannelInfo) error
	UpdateChannel(channel model.ChannelInfo) error
	DeleteChannel(id string) error
	AddTrackToSchedule(id, trackid string, start, end time.Time) error
	ChangeChannelStatus(id string, status model.ChannelStatus) error
	ChangeLogo(id, logo string) error
}

func NewChannelDB() ChannelDB {
	return &ChannelDBImpl{
		conn: database.GetConnection(),
	}
}

type ChannelDBImpl struct {
	conn *sql.DB
}

func (db *ChannelDBImpl) GetChannels() ([]model.ChannelShortInfo, error) {
	res := make([]model.ChannelShortInfo, 0)
	rows, err := db.conn.Query("SELECT id, title, logo FROM channels WHERE status=$1", model.ActiveChannel)
	if err != nil {
		return res, err
	}
	for rows.Next() {
		var temp model.ChannelShortInfo
		var logo sql.NullString
		err = rows.Scan(&temp.ID, &temp.Title, &logo)
		if err != nil {
			return res, err
		}
		if logo.Valid {
			temp.Logo = logo.String
		}
		res = append(res, temp)
	}

	return res, nil
}

func (db *ChannelDBImpl) GetChannelById(id string) (*model.ChannelInfo, error) {
	var res model.ChannelInfo
	rows, err := db.conn.Query("SELECT id, title, description, logo, status FROM channels WHERE id=$1", id)
	if err != nil {
		return &res, err
	}
	if rows.Next() {
		err = rows.Scan(&res.ID, &res.Title, &res.Description, &res.Logo, &res.Status)
		return &res, err
	}

	return nil, nil
}

func (db *ChannelDBImpl) GetCurrentTrack(id string) (string, error) {
	var res string
	rows, err := db.conn.Query("SELECT audio FROM schedule JOIN tracks ON tracks.id=schedule.trackid WHERE channelid=$1 AND NOW() between startdate AND enddate LIMIT 1", id)
	if err != nil {
		return res, err
	}
	if rows.Next() {
		err = rows.Scan(&res)
		return res, err
	}

	return res, errors.New("Track not found")

}
func (db *ChannelDBImpl) CreateChannel(channel model.ChannelInfo) error {
	_, err := db.conn.Exec("INSERT INTO channels (title, description, status) VALUES ($1, $2, $3)", channel.Title, channel.Description, channel.Status)
	return err
}

func (db *ChannelDBImpl) UpdateChannel(channel model.ChannelInfo) error {
	_, err := db.conn.Exec("UPDATE channels SET title=$1, description=$2, status=$3 WHERE id=$4", channel.Title, channel.Description, channel.Status, channel.ID)
	return err
}

func (db *ChannelDBImpl) DeleteChannel(id string) error {
	_, err := db.conn.Exec("DELETE FROM channels WHERE id=$1", id)
	return err
}

func (db *ChannelDBImpl) AddTrackToSchedule(id, trackid string, start, end time.Time) error {
	_, err := db.conn.Exec("INSERT INTO schedule (channelid, trackid, start, end)")
	return err
}

func (db *ChannelDBImpl) ChangeChannelStatus(id string, status model.ChannelStatus) error {
	_, err := db.conn.Exec("UPDATE channels SET status=$1 WHERE id=$2", status, id)
	return err
}

func (db *ChannelDBImpl) ChangeLogo(id, logo string) error {
	_, err := db.conn.Exec("UPDATE channels SET logo=$1 WHERE id=$2", logo, id)
	return err
}
