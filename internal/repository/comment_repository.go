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
	parent := &comment.Parent
	if comment.Parent <= 0 {
		parent = nil
	}
	err := db.conn.QueryRow("INSERT INTO comments (userid, parent, text, time) VALUES ($1, $2, $3, NOW()) RETURNING id", comment.UserID, parent, comment.Text).Scan(&id)
	return id, err
}
