package service

import (
	"fmt"
	"netradio/internal/controller/requests"
	"netradio/internal/controller/responses"
	"netradio/internal/model"
	"netradio/internal/repository"
	"netradio/pkg/email"
	"netradio/pkg/errors"
	"netradio/pkg/util"
	"strconv"
	"time"
)

type UserService interface {
	Register(r requests.RegisterRequest) (responses.RegisterResponse, int, error)
	Login(r requests.LoginRequest) (responses.LoginResponse, int, error)
	Refresh(r requests.RefreshRequest) (responses.RefreshResponse, int, error)
	Logout(r requests.LogoutRequest) error
	DeleteUser(r requests.DeleteUserRequset) error
	ResetPasswordSendCode(r requests.ResetPasswordSendCodeRequest) error
	ResetPasswordVerifyCode(r requests.ResetPasswordVerifyCodeRequest) (responses.ResetPasswordVerifyCodeResponse, int, error)
	ResetPasswordChange(r requests.ResetPasswordChangeRequest) (responses.ResetPasswordChangeResponse, error)
}

func NewUserService() UserService {
	return &UserServiceImpl{
		db: repository.NewUserDB(),
	}
}

type UserServiceImpl struct {
	db repository.UserDB
}

func (s *UserServiceImpl) Register(r requests.RegisterRequest) (responses.RegisterResponse, int, error) {
	var res responses.RegisterResponse
	if !util.IsEmailValid(r.Email) {
		res.Error = errors.InvalidEmail
		return res, 0, nil
	}
	if !util.IsPasswordValid(r.Password) {
		res.Error = errors.InvalidPassword
		return res, 0, nil
	}

	user, err := s.db.GetUserByEmail(r.Email)
	if err != nil {
		return res, 0, err
	}
	if user != nil {
		res.Error = errors.UserAlreadyExist
		return res, 0, nil
	}

	hashedPassword, err := util.HashPassword(r.Password)
	if err != nil {
		return res, 0, err
	}

	var newUser model.User
	newUser.Name = r.Name
	newUser.Email = r.Email
	newUser.Password = hashedPassword
	newUser.Role = model.UserRegistered
	err = s.db.CreateUser(newUser)
	if err != nil {
		return res, 0, err
	}

	user, err = s.db.GetUserByEmail(r.Email)
	if err != nil {
		return res, 0, err
	}

	return res, user.ID, err
}

func (s *UserServiceImpl) Login(r requests.LoginRequest) (responses.LoginResponse, int, error) {
	var res responses.LoginResponse
	user, err := s.db.GetUserByEmail(r.Email)
	if err != nil {
		return res, 0, err
	}
	if user == nil {
		res.Error = errors.UserNotFound
		return res, 0, err
	}
	if !util.CheckPassword(r.Password, user.Password) {
		res.Error = errors.WrongPassword
		return res, 0, err
	}

	return res, user.ID, err
}

func (s *UserServiceImpl) Refresh(r requests.RefreshRequest) (responses.RefreshResponse, int, error) {
	var res responses.RefreshResponse
	sessions, err := s.db.GetSessionsByRefreshToken(r.RefreshToken)
	if err != nil {
		return res, 0, err
	}

	var session *model.Session
	for i := range sessions {
		if sessions[i].IP == r.IP {
			session = &sessions[i]
		}
	}
	if session == nil {
		res.Error = errors.SessionNotFound
		return res, 0, nil
	}

	if session.Expires.Before(time.Now()) {
		err = s.db.DeleteSession(strconv.Itoa(session.UserID))
		if err != nil {
			return res, 0, err
		}
		res.Error = errors.SessionExpired
		return res, 0, nil
	}

	return res, session.UserID, err
}

func (s *UserServiceImpl) Logout(r requests.LogoutRequest) error {
	return s.db.DeleteSession(r.UserID)
}

func (s *UserServiceImpl) DeleteUser(r requests.DeleteUserRequset) error {
	return s.db.DeleteUser(r.UserID)
}

func (s *UserServiceImpl) ResetPasswordSendCode(r requests.ResetPasswordSendCodeRequest) error {
	user, err := s.db.GetUserByEmail(r.Email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("User with email " + r.Email + " not found")
	}

	code := util.GenVerificationCode(r.Email)
	err = s.db.CreateVerificationCode(code)
	if err != nil {
		return err
	}

	body := fmt.Sprintf("Your code for password reset is %s.", code.Value)
	go email.SendMail(user.Email, "Reset password", body)

	return nil
}

func (s *UserServiceImpl) ResetPasswordVerifyCode(r requests.ResetPasswordVerifyCodeRequest) (responses.ResetPasswordVerifyCodeResponse, int, error) {
	var res responses.ResetPasswordVerifyCodeResponse
	code, err := s.db.GetVerificationCodeByEmail(r.Email)
	if err != nil {
		return res, 0, err
	}

	if code == nil || code.Value != r.Code {
		res.Error = errors.WrongCode
		return res, 0, nil
	}

	err = s.db.DeleteVerificationCode(r.Email)
	if err != nil {
		return res, 0, err
	}

	if code.Expires.Before(time.Now()) {
		res.Error = errors.CodeExpires
		return res, 0, nil
	}

	user, err := s.db.GetUserByEmail(r.Email)
	if err != nil {
		return res, 0, err
	}
	if user == nil {
		res.Error = errors.UserNotFound
		return res, 0, nil
	}

	return res, user.ID, nil
}

func (s *UserServiceImpl) ResetPasswordChange(r requests.ResetPasswordChangeRequest) (responses.ResetPasswordChangeResponse, error) {
	var res responses.ResetPasswordChangeResponse
	if !util.IsPasswordValid(r.NewPassword) {
		res.Error = errors.InvalidPassword
		return res, nil
	}

	hashedPassword, err := util.HashPassword(r.NewPassword)
	if err != nil {
		return res, err
	}

	err = s.db.ChangePassword(r.UserID, hashedPassword)
	if err != nil {
		return res, err
	}

	return res, nil
}
