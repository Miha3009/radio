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
	CreateNews(news model.News) (int, error)
	UpdateNews(news model.News) error
	DeleteNews(id string) error
	ChangeImage(id, image string) error
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
	rows, err := db.conn.Query("SELECT id, title, publication_date, image FROM news WHERE title LIKE $3 ORDER BY publication_date OFFSET $1 LIMIT $2", offset, limit, query)
	if err != nil {
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		var temp model.NewsShortInfo
		var image sql.NullString
		err = rows.Scan(&temp.ID, &temp.Title, &temp.PublicationDate, &image)
		if err != nil {
			return res, err
		}
		if image.Valid {
			temp.Image = image.String
		}
		res = append(res, temp)
	}

	return res, nil
}

func (db *NewsDBImpl) GetNewsById(id string) (*model.News, error) {
	var res model.News
	rows, err := db.conn.Query("SELECT id, title, content, publication_date, image FROM news WHERE id=$1", id)
	if err != nil {
		return &res, err
	}
	defer rows.Close()
	if rows.Next() {
		var image sql.NullString
		err = rows.Scan(&res.ID, &res.Title, &res.Content, &res.PublicationDate, &image)
		if image.Valid {
			res.Image = image.String
		}
		return &res, err
	}

	return nil, nil
}

func (db *NewsDBImpl) CreateNews(news model.News) (int, error) {
	var id int
	err := db.conn.QueryRow("INSERT INTO news (title, content, publication_date) VALUES ($1, $2, $3) RETURNING id", news.Title, news.Content, news.PublicationDate).Scan(&id)
	return id, err
}

func (db *NewsDBImpl) UpdateNews(news model.News) error {
	_, err := db.conn.Exec("UPDATE news SET title=$1, content=$2, publication_date=$3 WHERE id=$4", news.Title, news.Content, news.PublicationDate, news.ID)
	return err
}

func (db *NewsDBImpl) DeleteNews(id string) error {
	_, err := db.conn.Exec("DELETE FROM news WHERE id=$1", id)
	return err
}

func (db *NewsDBImpl) ChangeImage(id, image string) error {
	_, err := db.conn.Exec("UPDATE news SET image=$1 WHERE id=$2", image, id)
	return err
}
