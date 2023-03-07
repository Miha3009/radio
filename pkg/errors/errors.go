package errors

import "errors"

var (
	UserAlreadyExist = "error.user.exist"
	UserNotFound     = "error.user.notFound"
	InvalidEmail     = "error.email.invalid"
	InvalidPassword  = "error.password.invalid"
	WrongPassword    = "error.password.wrong"
	SessionExpired   = "error.session.expired"
	SessionNotFound  = "error.session.notFound"
	WrongCode        = "error.code.wrong"
	CodeExpires      = "erorr.code.expires"
)

func New(text string) error {
	return errors.New(text)
}

func Wrap(err error, message string) error {
	return errors.New(message + ": " + err.Error())
}
