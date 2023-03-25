package repository

import (
	"database/sql"
	"fmt"
	"netradio/internal/model"
	"netradio/pkg/database"
	"time"
)

type ChannelDB interface {
	GetChannels(offset, limit int, query, status string) ([]model.ChannelShortInfo, int, error)
	GetChannelById(id string) (*model.ChannelInfo, error)
	GetCurrentTrack(id string) (*model.Track, error)
	CreateChannel(channel model.ChannelInfo) (int, error)
	UpdateChannel(channel model.ChannelInfo) error
	DeleteChannel(id string) error
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

func (db *ChannelDBImpl) GetChannels(offset, limit int, query, status string) ([]model.ChannelShortInfo, int, error) {
	res := make([]model.ChannelShortInfo, 0)
	statusString := ""
	if status == "active" {
		statusString = " AND status=1"
	} else if status == "stopped" {
		statusString = " AND status=0"
	}
	query = "%" + query + "%"
	rows, err := db.conn.Query(fmt.Sprintf("SELECT id, title, logo, status FROM channels WHERE title LIKE $3%s ORDER BY id OFFSET $1 LIMIT $2", statusString), offset, limit, query)
	if err != nil {
		return res, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var temp model.ChannelShortInfo
		var logo sql.NullString
		err = rows.Scan(&temp.ID, &temp.Title, &logo, &temp.Status)
		if err != nil {
			return res, 0, err
		}
		if logo.Valid {
			temp.Logo = logo.String
		}
		res = append(res, temp)
	}

	count := 0
	err = db.conn.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM channels WHERE title LIKE $1%s", statusString), query).Scan(&count)

	return res, count, err
}

func (db *ChannelDBImpl) GetChannelById(id string) (*model.ChannelInfo, error) {
	var res model.ChannelInfo
	rows, err := db.conn.Query("SELECT id, title, description, logo, status FROM channels WHERE id=$1", id)
	if err != nil {
		return &res, err
	}
	defer rows.Close()
	if rows.Next() {
		var logo sql.NullString
		err = rows.Scan(&res.ID, &res.Title, &res.Description, &logo, &res.Status)
		if err != nil {
			return nil, err
		}
		if logo.Valid {
			res.Logo = logo.String
		}

		return &res, nil
	}

	return nil, nil
}

func (db *ChannelDBImpl) GetCurrentTrack(id string) (*model.Track, error) {
	var res model.Track
	rows, err := db.conn.Query("SELECT tracks.id, tracks.title, tracks.performancer, tracks.year, tracks.audio, tracks.duration FROM schedule JOIN tracks ON tracks.id=schedule.trackid WHERE channelid=$1 AND NOW() BETWEEN startdate AND enddate LIMIT 1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var audio sql.NullString
		var duration sql.NullInt64
		err = rows.Scan(&res.ID, &res.Title, &res.Performancer, &res.Year, &audio, &duration)
		if err != nil {
			return nil, err
		}
		if audio.Valid {
			res.Audio = audio.String
		}
		if duration.Valid {
			res.Duration = time.Duration(duration.Int64)
		} else {
			res.Duration = time.Minute
		}

		return &res, nil
	}

	return nil, nil

}
func (db *ChannelDBImpl) CreateChannel(channel model.ChannelInfo) (int, error) {
	var id int
	err := db.conn.QueryRow("INSERT INTO channels (title, description, status) VALUES ($1, $2, $3) RETURNING id", channel.Title, channel.Description, channel.Status).Scan(&id)
	return id, err
}

func (db *ChannelDBImpl) UpdateChannel(channel model.ChannelInfo) error {
	_, err := db.conn.Exec("UPDATE channels SET title=$1, description=$2, status=$3 WHERE id=$4", channel.Title, channel.Description, channel.Status, channel.ID)
	return err
}

func (db *ChannelDBImpl) DeleteChannel(id string) error {
	_, err := db.conn.Exec("DELETE FROM channels WHERE id=$1", id)
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
