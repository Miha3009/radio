package repository

import (
	"database/sql"
	"netradio/internal/model"
	"netradio/pkg/database"
)

type CommentDB interface {
	CreateComment(comment model.Comment) (int, error)
}

func NewCommentDB() CommentDB {
	return &CommentDBImpl{
		conn: database.GetConnection(),
	}
}

type CommentDBImpl struct {
	conn *sql.DB
}

func (db *CommentDBImpl) CreateComment(comment model.Comment) (int, error) {
	id := 0
	err := db.conn.QueryRow("INSERT INTO comments (userid, parent, text, time) VALUES ($1, $2, $3, $4) RETURNING id", comment.UserID, comment.Parent, comment.Text, comment.Date).Scan(&id)
	return id, err
}
