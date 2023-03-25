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
	webrtc "netradio/pkg/webrtc"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/net/websocket"
)

func HandleGetChannels(ctx context.Context) (any, error) {
	var request requests.GetChannelsRequest
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
	request.Status = ctx.GetRequest().URL.Query().Get("status")

	res, err := service.NewChannelService().GetChannels(request)
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

	res, err := service.NewChannelService().CreateChannel(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return res, nil
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

func HandleUploadLogo(ctx context.Context) (any, error) {
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

	var request requests.UploadLogoRequest
	request.ID = chi.URLParam(ctx.GetRequest(), "id")
	request.Logo = path

	res, err := service.NewChannelService().UploadLogo(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return res, nil
}

func HandleAddTrack(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role != model.UserAdministrator {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.AddTrackToScheduleRequest
	request.ChannelID = chi.URLParam(ctx.GetRequest(), "id")
	decoder := json.NewDecoder(ctx.GetRequest().Body)
	err := decoder.Decode(&request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	err = service.NewChannelService().AddTrack(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}

func HandleGetCurrentTrack(ctx context.Context) (any, error) {
	var request requests.GetCurrentTrackRequest
	request.ID = chi.URLParam(ctx.GetRequest(), "id")

	user := ctx.GetUser()
	if user.Role != model.UserGuest {
		userId := strconv.Itoa(user.ID)
		request.UserID = &userId
	}

	res, err := service.NewChannelService().GetCurrentTrack(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return res, nil
}

func HandleSchedule(ctx context.Context) (any, error) {
	var request requests.GetScheduleRequest
	request.ChannelID = chi.URLParam(ctx.GetRequest(), "id")

	past, err := strconv.Atoi(ctx.GetRequest().URL.Query().Get("past"))
	if err != nil {
		past = 1000000
	}
	request.Past = past
	next, err := strconv.Atoi(ctx.GetRequest().URL.Query().Get("next"))
	if err != nil {
		next = 1000000
	}
	request.Next = next

	res, err := service.NewChannelService().GetSchedule(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return res, nil
}

func HandleScheduleRange(ctx context.Context) (any, error) {
	var request requests.GetScheduleRangeRequest
	request.ChannelID = chi.URLParam(ctx.GetRequest(), "id")

	from, err := time.Parse(time.RFC3339, ctx.GetRequest().URL.Query().Get("from"))
	if err != nil {
		from = time.Now().Add(-time.Hour * 1000000)
	}
	request.From = from
	to, err := time.Parse(time.RFC3339, ctx.GetRequest().URL.Query().Get("to"))
	if err != nil {
		to = time.Now().Add(time.Hour * 1000000)
	}
	request.To = to

	res, err := service.NewChannelService().GetScheduleRange(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return res, nil
}

func RouteChannelPaths(
	core handlers.Core,
	router chi.Router,
) {
	router.MethodFunc("GET", "/channel", handlers.MakeHandler(handlers.MakeJSONWrapper(HandleGetChannels), core))
	router.MethodFunc("GET", "/channel/{id}", handlers.MakeHandler(handlers.MakeJSONWrapper(HandleGetChannel), core))
	router.MethodFunc("PUT", "/channel", handlers.MakeHandler(handlers.MakeJSONWrapper(HandleCreateChannel), core))
	router.MethodFunc("DELETE", "/channel/{id}", handlers.MakeHandler(HandleDeleteChannel, core))
	router.MethodFunc("PATCH", "/channel/{id}", handlers.MakeHandler(HandleUpdateChannel, core))
	router.MethodFunc("POST", "/channel/{id}/start", handlers.MakeHandler(HandleStartChannel, core))
	router.MethodFunc("POST", "/channel/{id}/stop", handlers.MakeHandler(HandleStopChannel, core))
	router.MethodFunc("POST", "/channel/{id}/upload-logo", handlers.MakeHandler(handlers.MakeJSONWrapper(HandleUploadLogo), core))
	router.MethodFunc("POST", "/channel/{id}/add-track", handlers.MakeHandler(HandleAddTrack, core))
	router.MethodFunc("GET", "/channel/{id}/track", handlers.MakeHandler(handlers.MakeJSONWrapper(HandleGetCurrentTrack), core))
	router.MethodFunc("GET", "/channel/{id}/schedule", handlers.MakeHandler(handlers.MakeJSONWrapper(HandleSchedule), core))
	router.MethodFunc("GET", "/channel/{id}/schedule-range", handlers.MakeHandler(handlers.MakeJSONWrapper(HandleScheduleRange), core))

	router.HandleFunc("/channel/{id}/connect",
		func(w http.ResponseWriter, req *http.Request) {
			conn := webrtc.ConnectStruct{ID: chi.URLParam(req, "id")}
			s := websocket.Server{Handler: websocket.Handler(conn.HandleConnectChannel)}
			s.ServeHTTP(w, req)
		})
}
