package controller

import (
	"encoding/json"
	"net/http"
	"netradio/internal/controller/requests"
	"netradio/internal/model"
	"netradio/internal/service"
	"netradio/pkg/context"
	"netradio/pkg/handlers"

	"github.com/go-chi/chi/v5"
)

func HandleDeleteTrackFromSchedule(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role != model.UserAdministrator {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.DeleteTrackFromScheduleRequest
	request.ID = chi.URLParam(ctx.GetRequest(), "id")

	err := service.NewScheduleService().DeleteTrack(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}

func HandleUpdateTracksFromSchedule(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role != model.UserAdministrator {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.UpdateNewsRequest
	request.ID = chi.URLParam(ctx.GetRequest(), "id")
	decoder := json.NewDecoder(ctx.GetRequest().Body)
	err := decoder.Decode(&request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	res, err := service.NewNewsService().UpdateNews(request)
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

func RouteSchedulePaths(
	core handlers.Core,
	router chi.Router,
) {
	router.MethodFunc("DELETE", "/schedule/{id}", handlers.MakeHandler(HandleDeleteTrackFromSchedule, core))
	router.MethodFunc("PATCH", "/schedule", handlers.MakeHandler(HandleUpdateTracksFromSchedule, core))
}
