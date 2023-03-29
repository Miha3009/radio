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
	Register(r requests.RegisterRequest) (responses.RegisterResponse, error)
	Login(r requests.LoginRequest) (responses.LoginResponse, error)
	Refresh(r requests.RefreshRequest) (responses.RefreshResponse, error)
	Logout(r requests.LogoutRequest) error
	DeleteUser(r requests.DeleteUserRequset) error
	ResetPasswordSendCode(r requests.ResetPasswordSendCodeRequest) error
	ResetPasswordVerifyCode(r requests.ResetPasswordVerifyCodeRequest) (responses.ResetPasswordVerifyCodeResponse, error)
	ResetPasswordChange(r requests.ResetPasswordChangeRequest) (responses.ResetPasswordChangeResponse, error)
	UploadAvatar(r requests.UploadAvatarRequest) error
	UpdateUser(r requests.UpdateUserRequest) error
}

func NewUserService() UserService {
	return &UserServiceImpl{
		db: repository.NewUserDB(),
	}
}

type UserServiceImpl struct {
	db repository.UserDB
}

func (s *UserServiceImpl) Register(r requests.RegisterRequest) (responses.RegisterResponse, error) {
	var res responses.RegisterResponse
	if !util.IsEmailValid(r.Email) {
		res.Error = errors.InvalidEmail
		return res, nil
	}
	if !util.IsPasswordValid(r.Password) {
		res.Error = errors.InvalidPassword
		return res, nil
	}

	user, err := s.db.GetUserByEmail(r.Email)
	if err != nil {
		return res, err
	}
	if user != nil {
		res.Error = errors.UserAlreadyExist
		return res, nil
	}

	hashedPassword, err := util.HashPassword(r.Password)
	if err != nil {
		return res, err
	}

	var newUser model.User
	newUser.Name = r.Name
	newUser.Email = r.Email
	newUser.Password = hashedPassword
	newUser.Role = model.UserRegistered
	err = s.db.CreateUser(newUser)
	if err != nil {
		return res, err
	}

	user, err = s.db.GetUserByEmail(r.Email)
	if err != nil {
		return res, err
	}
	res.UserID = user.ID

	return res, err
}

func (s *UserServiceImpl) Login(r requests.LoginRequest) (responses.LoginResponse, error) {
	var res responses.LoginResponse
	user, err := s.db.GetUserByEmail(r.Email)
	if err != nil {
		return res, err
	}
	if user == nil {
		res.Error = errors.UserNotFound
		return res, err
	}
	if !util.CheckPassword(r.Password, user.Password) {
		res.Error = errors.WrongPassword
		return res, err
	}
	res.UserID = user.ID

	return res, err
}

func (s *UserServiceImpl) Refresh(r requests.RefreshRequest) (responses.RefreshResponse, error) {
	var res responses.RefreshResponse
	sessions, err := s.db.GetSessionsByRefreshToken(r.RefreshToken)
	if err != nil {
		return res, err
	}

	var session *model.Session
	for i := range sessions {
		if sessions[i].IP == r.IP {
			session = &sessions[i]
		}
	}
	if session == nil {
		res.Error = errors.SessionNotFound
		return res, nil
	}

	if session.Expires.Before(time.Now()) {
		err = s.db.DeleteSession(strconv.Itoa(session.UserID))
		if err != nil {
			return res, err
		}
		res.Error = errors.SessionExpired
		return res, nil
	}
	res.UserID = session.UserID

	return res, err
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

func (s *UserServiceImpl) ResetPasswordVerifyCode(r requests.ResetPasswordVerifyCodeRequest) (responses.ResetPasswordVerifyCodeResponse, error) {
	var res responses.ResetPasswordVerifyCodeResponse
	code, err := s.db.GetVerificationCodeByEmail(r.Email)
	if err != nil {
		return res, err
	}

	if code == nil || code.Value != r.Code {
		res.Error = errors.WrongCode
		return res, nil
	}

	err = s.db.DeleteVerificationCode(r.Email)
	if err != nil {
		return res, err
	}

	if code.Expires.Before(time.Now()) {
		res.Error = errors.CodeExpires
		return res, nil
	}

	user, err := s.db.GetUserByEmail(r.Email)
	if err != nil {
		return res, err
	}
	if user == nil {
		res.Error = errors.UserNotFound
		return res, nil
	}
	res.UserID = user.ID

	return res, nil
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

func (s *UserServiceImpl) UploadAvatar(r requests.UploadAvatarRequest) error {
	return s.db.ChangeAvatar(r.UserID, r.Avatar)
}

func (s *UserServiceImpl) UpdateUser(r requests.UpdateUserRequest) error {
	if r.Email != nil {
		r.User.Email = *r.Email
	}
	if r.Name != nil {
		r.User.Name = *r.Name
	}

	return s.db.UpdateUser(r.User)
}
