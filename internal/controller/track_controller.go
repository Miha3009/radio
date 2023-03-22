package controller

import (
	"encoding/json"
	"net/http"
	"netradio/internal/controller/requests"
	"netradio/internal/model"
	"netradio/internal/service"
	"netradio/pkg/context"
	"netradio/pkg/files"
	"netradio/pkg/handlers"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func HandleGetTrack(ctx context.Context) (any, error) {
	var request requests.GetTrackRequest
	request.ID = chi.URLParam(ctx.GetRequest(), "id")

	user := ctx.GetUser()
	if user.Role != model.UserGuest {
		userId := strconv.Itoa(user.ID)
		request.UserID = &userId
	}

	res, err := service.NewTrackService().GetTrack(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}
	if !res.Found {
		ctx.GetResponseWriter().WriteHeader(http.StatusNotFound)
		return res, nil
	}

	return res, nil
}

func HandleGetTrackList(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role != model.UserAdministrator {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.GetTrackListRequest
	offset, err := strconv.Atoi(ctx.GetRequest().URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}
	request.Offset = offset
	limit, err := strconv.Atoi(ctx.GetRequest().URL.Query().Get("limit"))
	if err != nil {
		limit = 1000000
	}
	request.Limit = limit
	request.Query = ctx.GetRequest().URL.Query().Get("query")

	res, err := service.NewTrackService().GetTrackList(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return res, nil
}

func HandleCreateTrack(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role != model.UserAdministrator {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.CreateTrackRequest
	decoder := json.NewDecoder(ctx.GetRequest().Body)
	err := decoder.Decode(&request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	err = service.NewTrackService().CreateTrack(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}

func HandleDeleteTrack(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role != model.UserAdministrator {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.DeleteTrackRequest
	request.ID = chi.URLParam(ctx.GetRequest(), "id")

	err := service.NewTrackService().DeleteTrack(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}

func HandleUpdateTrack(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role != model.UserAdministrator {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.UpdateTrackRequest
	request.ID = chi.URLParam(ctx.GetRequest(), "id")
	decoder := json.NewDecoder(ctx.GetRequest().Body)
	err := decoder.Decode(&request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	res, err := service.NewTrackService().UpdateTrack(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}
	if !res.Found {
		ctx.GetResponseWriter().WriteHeader(http.StatusNotFound)
		return res, nil
	}

	return nil, nil
}

func HandleLikeTrack(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role == model.UserGuest {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.LikeTrackRequest
	decoder := json.NewDecoder(ctx.GetRequest().Body)
	err := decoder.Decode(&request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusBadRequest)
		return nil, err
	}
	request.ID = chi.URLParam(ctx.GetRequest(), "id")
	request.UserID = user.ID

	err = service.NewTrackService().LikeTrack(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}

func HandleCommentTrack(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role == model.UserGuest {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.CommentTrackRequest
	decoder := json.NewDecoder(ctx.GetRequest().Body)
	err := decoder.Decode(&request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusBadRequest)
		return nil, err
	}
	request.ID = chi.URLParam(ctx.GetRequest(), "id")
	request.UserID = user.ID

	err = service.NewTrackService().CommentTrack(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}

func HandleUploadTrack(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role != model.UserAdministrator {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	path, err := files.Save(ctx.GetRequest())
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	var request requests.UploadTrackRequest
	request.ID = chi.URLParam(ctx.GetRequest(), "id")
	request.Audio = path

	err = service.NewTrackService().UploadTrack(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}

func RouteTrackPaths(
	core handlers.Core,
	router chi.Router,
) {
	router.MethodFunc("GET", "/track", handlers.MakeHandler(handlers.MakeJSONWrapper(HandleGetTrackList), core))
	router.MethodFunc("GET", "/track/{id}", handlers.MakeHandler(handlers.MakeJSONWrapper(HandleGetTrack), core))
	router.MethodFunc("PUT", "/track", handlers.MakeHandler(HandleCreateTrack, core))
	router.MethodFunc("DELETE", "/track/{id}", handlers.MakeHandler(HandleDeleteTrack, core))
	router.MethodFunc("PATCH", "/track/{id}", handlers.MakeHandler(HandleUpdateTrack, core))
	router.MethodFunc("POST", "/track/{id}/like", handlers.MakeHandler(HandleLikeTrack, core))
	router.MethodFunc("POST", "/track/{id}/comment", handlers.MakeHandler(HandleCommentTrack, core))
	router.MethodFunc("POST", "/track/{id}/upload", handlers.MakeHandler(HandleUploadTrack, core))
}
