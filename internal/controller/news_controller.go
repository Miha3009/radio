package controller

import (
	"encoding/json"
	"net/http"
	"netradio/internal/controller/requests"
	"netradio/internal/service"
	"netradio/pkg/context"
	"netradio/pkg/handlers"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func HandleListNews(ctx context.Context) (any, error) {
	var request requests.ListNewsRequest
	decoder := json.NewDecoder(ctx.GetRequest().Body)
	err := decoder.Decode(&request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	news, err := service.NewNewsService().GetRange(request.Offset, request.Limit)

	if err != nil {
		return nil, err
	}

	return news, nil
}

func HandleGetNews(ctx context.Context) (any, error) {
	idParam := chi.URLParam(ctx.GetRequest(), "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	news, err := service.NewNewsService().Get(id)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusNotFound)
		return nil, err
	}

	return news, nil
}

func RouteNewsPaths(
	core handlers.Core,
	router chi.Router,
) {
	router.MethodFunc("GET", "/news", handlers.MakeHandler(handlers.MakeJSONWrapper(HandleListNews), core))
	router.MethodFunc("GET", "/news/{id}", handlers.MakeHandler(handlers.MakeJSONWrapper(HandleGetNews), core))
}
