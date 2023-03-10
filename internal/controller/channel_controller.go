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

func HandleGetChannels(ctx context.Context) (any, error) {
	res, err := service.NewChannelService().GetChannels()
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return res, nil
}

func HandleGetChannel(ctx context.Context) (any, error) {
	var request requests.GetChannelRequest
	request.ID = chi.URLParam(ctx.GetRequest(), "id")

	res, err := service.NewChannelService().GetChannel(request)
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

func HandleCreateChannel(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role != model.UserAdministrator {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.CreateChannelRequest
	decoder := json.NewDecoder(ctx.GetRequest().Body)
	err := decoder.Decode(&request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	err = service.NewChannelService().CreateChannel(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}

func HandleDeleteChannel(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role != model.UserAdministrator {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.DeleteChannelRequest
	request.ID = chi.URLParam(ctx.GetRequest(), "id")

	err := service.NewChannelService().DeleteChannel(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}

func HandleUpdateChannel(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role != model.UserAdministrator {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.UpdateChannelRequest
	request.ID = chi.URLParam(ctx.GetRequest(), "id")
	decoder := json.NewDecoder(ctx.GetRequest().Body)
	err := decoder.Decode(&request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	res, err := service.NewChannelService().UpdateChannel(request)
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

func HandleStartChannel(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role != model.UserAdministrator {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.StartChannelRequest
	request.ID = chi.URLParam(ctx.GetRequest(), "id")

	err := service.NewChannelService().StartChannel(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}

func HandleStopChannel(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role != model.UserAdministrator {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.StopChannelRequest
	request.ID = chi.URLParam(ctx.GetRequest(), "id")

	err := service.NewChannelService().StopChannel(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}

func RouteChannelPaths(
	core handlers.Core,
	router chi.Router,
) {
	router.MethodFunc("GET", "/channel", handlers.MakeHandler(handlers.MakeJSONWrapper(HandleGetChannels), core))
	router.MethodFunc("GET", "/channel/{id}", handlers.MakeHandler(handlers.MakeJSONWrapper(HandleGetChannel), core))
	router.MethodFunc("PUT", "/channel/create", handlers.MakeHandler(HandleCreateChannel, core))
	router.MethodFunc("DELETE", "/channel/{id}", handlers.MakeHandler(HandleDeleteChannel, core))
	router.MethodFunc("PATCH", "/channel/{id}", handlers.MakeHandler(HandleUpdateChannel, core))
	router.MethodFunc("POST", "/channel/{id}/start", handlers.MakeHandler(HandleStartChannel, core))
	router.MethodFunc("POST", "/channel/{id}/stop", handlers.MakeHandler(HandleStopChannel, core))
}
