package repository

import (
	"database/sql"
	"errors"
	"netradio/internal/model"
	"netradio/pkg/database"
	"strconv"
)

type NewsDB interface {
	Count() (int, error)
	GetRange(offset, limit int) ([]model.News, error)
	Get(id int) (model.News, error)
	Add(news model.News) int
}

func NewNewsDB() NewsDB {
	return &NewsDBImpl{
		conn: database.GetConnection(),
	}
}

type NewsDBImpl struct {
	conn *sql.DB
}

func (db *NewsDBImpl) Count() (int, error) {
	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM news").Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (db *NewsDBImpl) GetRange(offset, limit int) ([]model.News, error) {
	rows, err := db.conn.Query("SELECT id, title, content FROM news ORDER BY time OFFSET $1 LIMIT $2", offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]model.News, 0)
	for rows.Next() {
		var row model.News
		if err := rows.Scan(&row.ID, &row.Title, &row.Content); err != nil {
			return nil, err
		}
		res = append(res, row)
	}
	return res, nil
}

func (db *NewsDBImpl) Get(id int) (model.News, error) {
	var res model.News
	err := db.conn.QueryRow("SELECT id, title, content FROM news WHERE id=$1", id).Scan(&res.ID, &res.Title, &res.Content)
	if err != nil {
		return res, errors.New("no news for id=" + strconv.Itoa(id))
	}

	return res, nil
}

func (ds *NewsDBImpl) Add(model model.News) int {
	return model.ID
}
