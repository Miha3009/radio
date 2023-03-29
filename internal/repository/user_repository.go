package repository

import (
	"database/sql"
	"netradio/internal/model"
	"netradio/pkg/database"
	"netradio/pkg/files"
	"time"
)

type UserDB interface {
	GetUserByEmail(email string) (*model.User, error)
	GetUserById(id string) (*model.User, error)
	GetSessionsByRefreshToken(refreshToken string) ([]model.Session, error)
	GetVerificationCodeByEmail(email string) (*model.VerificationCode, error)
	CreateUser(user model.User) error
	CreateSession(userID int, refreshToken string, expires time.Time, ip string) error
	CreateVerificationCode(code model.VerificationCode) error
	UpdateUser(user model.User) error
	ChangeAvatar(id, newAvatar string) error
	ChangePassword(id, newPassword string) error
	DeleteUser(id string) error
	DeleteSession(userId string) error
	DeleteVerificationCode(email string) error
}

func NewUserDB() UserDB {
	return &UserDBImpl{
		conn: database.GetConnection(),
	}
}

type UserDBImpl struct {
	conn *sql.DB
}

func (db *UserDBImpl) GetUserByEmail(email string) (*model.User, error) {
	var res model.User
	rows, err := db.conn.Query("SELECT id, email, password, name, avatar, role FROM users WHERE email=$1", email)
	if err != nil {
		return &res, err
	}
	defer rows.Close()
	if rows.Next() {
		var avatar sql.NullString
		err = rows.Scan(&res.ID, &res.Email, &res.Password, &res.Name, &avatar, &res.Role)
		if avatar.Valid {
			res.Avatar = files.ToURL(avatar.String)
		}
		return &res, err
	}

	return nil, nil
}

func (db *UserDBImpl) GetUserById(id string) (*model.User, error) {
	var res model.User
	rows, err := db.conn.Query("SELECT id, email, password, name, avatar, role FROM users WHERE id=$1", id)
	if err != nil {
		return &res, err
	}
	defer rows.Close()
	if rows.Next() {
		var avatar sql.NullString
		err = rows.Scan(&res.ID, &res.Email, &res.Password, &res.Name, &avatar, &res.Role)
		if avatar.Valid {
			res.Avatar = files.ToURL(avatar.String)
		}
		return &res, err
	}

	return nil, nil
}

func (db *UserDBImpl) GetSessionsByRefreshToken(refreshToken string) ([]model.Session, error) {
	res := make([]model.Session, 0)
	rows, err := db.conn.Query("SELECT userid, refresh_token, expires, ip FROM sessions WHERE refresh_token=$1", refreshToken)
	if err != nil {
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		var temp model.Session
		err = rows.Scan(&temp.UserID, &temp.RefreshToken, &temp.Expires, &temp.IP)
		if err != nil {
			return res, err
		}
		res = append(res, temp)
	}

	return res, nil
}

func (db *UserDBImpl) GetVerificationCodeByEmail(email string) (*model.VerificationCode, error) {
	var res model.VerificationCode
	rows, err := db.conn.Query("SELECT email, value, expires FROM verification_codes WHERE email=$1", email)
	if err != nil {
		return &res, err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&res.Email, &res.Value, &res.Expires)
		return &res, err
	}

	return nil, nil
}

func (db *UserDBImpl) CreateUser(user model.User) error {
	_, err := db.conn.Exec("INSERT INTO users (email, password, name, avatar, role) VALUES ($1, $2, $3, $4, $5)", user.Email, user.Password, user.Name, user.Avatar, user.Role)
	return err
}

func (db *UserDBImpl) CreateSession(userID int, refreshToken string, expires time.Time, ip string) error {
	_, err := db.conn.Exec("INSERT INTO sessions (userid, refresh_token, expires, ip) VALUES ($1, $2, $3, $4) ON CONFLICT (userid, ip) DO UPDATE SET refresh_token=EXCLUDED.refresh_token, expires=EXCLUDED.expires", userID, refreshToken, expires, ip)
	return err
}

func (db *UserDBImpl) CreateVerificationCode(code model.VerificationCode) error {
	_, err := db.conn.Exec("INSERT INTO verification_codes (email, value, expires) VALUES ($1, $2, $3) ON CONFLICT (email) DO UPDATE SET value=EXCLUDED.value, expires=EXCLUDED.expires", code.Email, code.Value, code.Expires)
	return err
}

func (db *UserDBImpl) UpdateUser(user model.User) error {
	_, err := db.conn.Exec("UPDATE users SET name=$1, email=$2 WHERE id=$3", user.Name, user.Email, user.ID)
	return err
}

func (db *UserDBImpl) ChangeAvatar(id, avatar string) error {
	_, err := db.conn.Exec("UPDATE users SET avatar=$1 WHERE id=$2", avatar, id)
	return err
}

func (db *UserDBImpl) ChangePassword(id, newPassword string) error {
	_, err := db.conn.Exec("UPDATE users SET password=$1 WHERE id=$2", newPassword, id)
	return err
}

func (db *UserDBImpl) DeleteUser(id string) error {
	_, err := db.conn.Exec("DELETE FROM users WHERE id=$1", id)
	return err
}

func (db *UserDBImpl) DeleteSession(userID string) error {
	_, err := db.conn.Exec("DELETE FROM sessions WHERE userid=$1", userID)
	return err
}

func (db *UserDBImpl) DeleteVerificationCode(email string) error {
	_, err := db.conn.Exec("DELETE FROM verification_codes WHERE email=$1", email)
	return err
}
