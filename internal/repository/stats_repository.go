package repository

import (
	"database/sql"
	"netradio/internal/model"
	"netradio/pkg/database"
	"time"
)

type StatsDB interface {
	AddListenerTimestamp(channelid string, count int) error
	GetListeners(channelid string, from, to time.Time) ([]model.ListenerStat, error)
}

func NewStatsDB() StatsDB {
	return &StatsDBImpl{
		conn: database.GetConnection(),
	}
}

type StatsDBImpl struct {
	conn *sql.DB
}

func (db *StatsDBImpl) AddListenerTimestamp(channelid string, count int) error {
	_, err := db.conn.Exec("INSERT INTO listeners_stat (channelid, time, count) VALUES ($1, NOW(), $2)", channelid, count)
	return err
}

func (db *StatsDBImpl) GetListeners(channelid string, from, to time.Time) ([]model.ListenerStat, error) {
	rows, err := db.conn.Query("SELECT time, count FROM listeners_stat WHERE channelid=$1 AND time BETWEEN $2 AND $3 ORDER BY time ASC", channelid, from, to)
	if err != nil {
		return nil, err
	}
	res := make([]model.ListenerStat, 0)
	defer rows.Close()
	for rows.Next() {
		var temp model.ListenerStat
		err = rows.Scan(&temp.Time, &temp.Count)
		if err != nil {
			return res, err
		}
		res = append(res, temp)
	}

	return res, nil
}
