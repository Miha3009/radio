package repository

import (
	"database/sql"
	"netradio/internal/model"
	"netradio/pkg/database"
	"time"
)

type NewsDB interface {
	GetNewsList(offset, limit int, query string) ([]model.NewsShortInfo, int, error)
	GetNewsById(id string) (*model.News, error)
	CreateNews(news model.News) (int, error)
	UpdateNews(news model.News) error
	DeleteNews(id string) error
	ChangeImage(id, image string) error
	IsNewsLiked(id, userId string) (bool, error)
	LikeNews(id, userId string) error
	UnlikeNews(id, userId string) error
	GetNewsComments(id string) ([]model.Comment, error)
	CommentNews(id, commentId string) error
	LikeCount(id string) (int, error)
}

func NewNewsDB() NewsDB {
	return &NewsDBImpl{
		conn: database.GetConnection(),
	}
}

type NewsDBImpl struct {
	conn *sql.DB
}

func (db *NewsDBImpl) GetNewsList(offset, limit int, query string) ([]model.NewsShortInfo, int, error) {
	res := make([]model.NewsShortInfo, 0)
	query = "%" + query + "%"
	rows, err := db.conn.Query("SELECT id, title, publication_date, image FROM news WHERE title LIKE $3 ORDER BY publication_date DESC OFFSET $1 LIMIT $2", offset, limit, query)
	if err != nil {
		return res, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var temp model.NewsShortInfo
		var image sql.NullString
		err = rows.Scan(&temp.ID, &temp.Title, &temp.PublicationDate, &image)
		if err != nil {
			return res, 0, err
		}
		if image.Valid {
			temp.Image = image.String
		}
		res = append(res, temp)
	}

	count := 0
	err = db.conn.QueryRow("SELECT COUNT(*) FROM news WHERE title LIKE $1", query).Scan(&count)

	return res, count, err
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

func (db *NewsDBImpl) IsNewsLiked(id, userId string) (bool, error) {
	rows, err := db.conn.Query("SELECT * FROM news_likes WHERE newsid=$1 AND userid=$2", id, userId)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (db *NewsDBImpl) LikeNews(id, userId string) error {
	_, err := db.conn.Exec("INSERT INTO news_likes (newsid, userid, time) VALUES ($1, $2, $3)", id, userId, time.Now())
	return err
}

func (db *NewsDBImpl) UnlikeNews(id, userId string) error {
	_, err := db.conn.Exec("DELETE FROM news_likes WHERE newsid=$1 AND userid=$2", id, userId)
	return err
}

func (db *NewsDBImpl) GetNewsComments(id string) ([]model.Comment, error) {
	res := make([]model.Comment, 0)
	rows, err := db.conn.Query("SELECT comments.id, users.id, users.name, users.avatar, comments.parent, comments.text, comments.time FROM news_comments JOIN comments ON news_comments.commentid=comments.id JOIN users ON comments.userid=users.id WHERE news_comments.newsid=$1", id)
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
			temp.UserAvatar = avatar.String
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

func (db *NewsDBImpl) CommentNews(id, commentId string) error {
	_, err := db.conn.Exec("INSERT INTO news_comments (newsid, commentid) VALUES ($1, $2)", id, commentId)
	return err
}

func (db *NewsDBImpl) LikeCount(id string) (int, error) {
	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM news_likes WHERE newsid=$1", id).Scan(&count)
	return count, err
}
