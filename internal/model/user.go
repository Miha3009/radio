package model

import "time"

type User struct {
	ID       int      `json:"id"`
	Email    string   `json:"email"`
	Password string   `json:"-"`
	Name     string   `json:"name"`
	Avatar   string   `json:"photo"`
	Role     UserRole `json:"status"`
}

type UserRole int

const (
	UserGuest         UserRole = 0
	UserRegistered             = 1
	UserAdministrator          = 2
)

type Session struct {
	UserID       int
	RefreshToken string
	IP           string
	Expires      time.Time
}

type VerificationCode struct {
	Email   string
	Value   string
	Expires time.Time
}
