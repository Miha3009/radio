package repository

import (
	"database/sql"
	"netradio/internal/model"
	"netradio/pkg/database"
)

type NewsDB interface {
	GetNewsCount() (int, error)
	GetNewsList(offset, limit int, query string) ([]model.NewsShortInfo, error)
	GetNewsById(id string) (*model.News, error)
	CreateNews(news model.News) error
	UpdateNews(news model.News) error
	DeleteNews(id string) error
}

func NewNewsDB() NewsDB {
	return &NewsDBImpl{
		conn: database.GetConnection(),
	}
}

type NewsDBImpl struct {
	conn *sql.DB
}

func (db *NewsDBImpl) GetNewsCount() (int, error) {
	count := 0
	err := db.conn.QueryRow("SELECT COUNT(*) FROM news").Scan(&count)
	return count, err
}

func (db *NewsDBImpl) GetNewsList(offset, limit int, query string) ([]model.NewsShortInfo, error) {
	res := make([]model.NewsShortInfo, 0)
	query = "%" + query + "%"
	rows, err := db.conn.Query("SELECT id, title, publication_date FROM news WHERE title LIKE $3 ORDER BY publication_date OFFSET $1 LIMIT $2", offset, limit, query)
	defer rows.Close()
	if err != nil {
		return res, err
	}
	for rows.Next() {
		var temp model.NewsShortInfo
		err = rows.Scan(&temp.ID, &temp.Title, &temp.PublicationDate)
		if err != nil {
			return res, err
		}
		res = append(res, temp)
	}

	return res, nil
}

func (db *NewsDBImpl) GetNewsById(id string) (*model.News, error) {
	var res model.News
	rows, err := db.conn.Query("SELECT id, title, content, pubication_date FROM news WHERE id=$1", id)
	defer rows.Close()
	if err != nil {
		return &res, err
	}
	if rows.Next() {
		err = rows.Scan(&res.ID, &res.Title, &res.Content, &res.PublicationDate)
		return &res, err
	}

	return nil, nil
}

func (db *NewsDBImpl) CreateNews(news model.News) error {
	_, err := db.conn.Exec("INSERT INTO news (title, content, publication_date)  VALUES ($1, $2, $3)", news.Title, news.Content, news.PublicationDate)
	return err
}

func (db *NewsDBImpl) UpdateNews(news model.News) error {
	_, err := db.conn.Exec("UPDATE news SET title=$1, content=$2, publication_date=$3 WHERE id=$4", news.Title, news.Content, news.PublicationDate, news.ID)
	return err
}

func (db *NewsDBImpl) DeleteNews(id string) error {
	_, err := db.conn.Exec("DELETE FROM news WHERE id=$1", id)
	return err
}
