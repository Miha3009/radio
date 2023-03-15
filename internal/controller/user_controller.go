package controller

import (
	"encoding/json"
	"net"
	"net/http"
	"netradio/internal/controller/requests"
	"netradio/internal/model"
	"netradio/internal/service"
	"netradio/pkg/context"
	"netradio/pkg/files"
	"netradio/pkg/handlers"
	"netradio/pkg/jwt"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func HandleRegister(ctx context.Context) (any, error) {
	var request requests.RegisterRequest
	decoder := json.NewDecoder(ctx.GetRequest().Body)
	err := decoder.Decode(&request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	res, err := service.NewUserService().Register(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}
	if res.Error == "" {
		accessToken, err := jwt.NewAccessToken(res.UserID)
		if err != nil {
			ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
			return nil, err
		}
		res.AccessToken = accessToken
		jwt.AddRefreshTokenToCookie(ctx.GetResponseWriter(), ctx.GetRequest(), res.UserID)
	}

	return res, nil
}

func HandleLogin(ctx context.Context) (any, error) {
	var request requests.LoginRequest
	decoder := json.NewDecoder(ctx.GetRequest().Body)
	err := decoder.Decode(&request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	res, err := service.NewUserService().Login(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}
	if res.Error == "" {
		accessToken, err := jwt.NewAccessToken(res.UserID)
		if err != nil {
			ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
			return nil, err
		}
		res.AccessToken = accessToken
		jwt.AddRefreshTokenToCookie(ctx.GetResponseWriter(), ctx.GetRequest(), res.UserID)
	}

	return res, nil
}

func HandleRefresh(ctx context.Context) (any, error) {
	var request requests.RefreshRequest
	for _, cookie := range ctx.GetRequest().Cookies() {
		if cookie.Name == "refreshToken" {
			request.RefreshToken = cookie.Value
			break
		}
	}

	if request.RefreshToken == "" {
		ctx.GetResponseWriter().WriteHeader(http.StatusBadRequest)
		return nil, nil
	}

	ip, _, err := net.SplitHostPort(ctx.GetRequest().RemoteAddr)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusBadRequest)
		return nil, nil
	}
	request.IP = ip

	res, err := service.NewUserService().Refresh(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return res, err
	}
	if res.Error == "" {
		accessToken, err := jwt.NewAccessToken(res.UserID)
		if err != nil {
			ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
			return nil, err
		}
		res.AccessToken = accessToken
		jwt.AddRefreshTokenToCookie(ctx.GetResponseWriter(), ctx.GetRequest(), res.UserID)
	}

	return res, nil
}

func HandleLogout(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role == model.UserGuest {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.LogoutRequest
	request.UserID = strconv.Itoa(user.ID)
	err := service.NewUserService().Logout(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	http.SetCookie(ctx.GetResponseWriter(), &http.Cookie{Name: "refreshToken", Value: "", Expires: time.Unix(0, 0)})

	return nil, nil
}

func HandleUserDelete(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role == model.UserGuest {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.DeleteUserRequset
	request.UserID = strconv.Itoa(user.ID)
	err := service.NewUserService().DeleteUser(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	http.SetCookie(ctx.GetResponseWriter(), &http.Cookie{Name: "refreshToken", Value: "", Expires: time.Unix(0, 0)})

	return nil, nil
}

func HandleResetPasswordSendCode(ctx context.Context) (any, error) {
	var request requests.ResetPasswordSendCodeRequest
	decoder := json.NewDecoder(ctx.GetRequest().Body)
	err := decoder.Decode(&request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	err = service.NewUserService().ResetPasswordSendCode(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}

func HandleResetPasswordVerifyCode(ctx context.Context) (any, error) {
	var request requests.ResetPasswordVerifyCodeRequest
	decoder := json.NewDecoder(ctx.GetRequest().Body)
	err := decoder.Decode(&request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	res, err := service.NewUserService().ResetPasswordVerifyCode(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}
	if res.Error == "" {
		accessToken, err := jwt.NewAccessToken(res.UserID)
		if err != nil {
			ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
			return nil, err
		}
		res.AccessToken = accessToken
		jwt.AddRefreshTokenToCookie(ctx.GetResponseWriter(), ctx.GetRequest(), res.UserID)
	}

	return res, nil
}

func HandleResetPasswordChange(ctx context.Context) (any, error) {
	var request requests.ResetPasswordChangeRequest
	decoder := json.NewDecoder(ctx.GetRequest().Body)
	err := decoder.Decode(&request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	user := ctx.GetUser()
	if user.Role == model.UserGuest {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}
	request.UserID = strconv.Itoa(user.ID)

	res, err := service.NewUserService().ResetPasswordChange(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return res, nil
}

func HandleUploadAvatar(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role == model.UserGuest {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	path, err := files.Save(ctx.GetRequest())
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	var request requests.UploadAvatarRequest
	request.UserID = strconv.Itoa(user.ID)
	request.Avatar = path

	err = service.NewUserService().UploadAvatar(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}

func RouteUserPaths(
	core handlers.Core,
	router chi.Router,
) {
	router.MethodFunc("POST", "/register", handlers.MakeHandler(handlers.MakeJSONWrapper(HandleRegister), core))
	router.MethodFunc("POST", "/login", handlers.MakeHandler(handlers.MakeJSONWrapper(HandleLogin), core))
	router.MethodFunc("GET", "/refresh", handlers.MakeHandler(handlers.MakeJSONWrapper(HandleRefresh), core))
	router.MethodFunc("POST", "/logout", handlers.MakeHandler(HandleLogout, core))
	router.MethodFunc("DELETE", "/user", handlers.MakeHandler(HandleUserDelete, core))
	router.MethodFunc("POST", "/reset-password/send-code", handlers.MakeHandler(HandleResetPasswordSendCode, core))
	router.MethodFunc("GET", "/reset-password/verify-code", handlers.MakeHandler(handlers.MakeJSONWrapper(HandleResetPasswordVerifyCode), core))
	router.MethodFunc("POST", "/reset-password/change", handlers.MakeHandler(HandleResetPasswordChange, core))
	router.MethodFunc("POST", "/upload-avatar", handlers.MakeHandler(HandleUploadAvatar, core))
}
