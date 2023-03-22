package controller

import (
	"encoding/json"
	"net/http"
	"netradio/internal/controller/requests"
	"netradio/internal/model"
	"netradio/internal/service"
	"netradio/pkg/context"
	"netradio/pkg/handlers"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func HandleGetNewsList(ctx context.Context) (any, error) {
	var request requests.GetNewsListRequest
	offset, err := strconv.Atoi(chi.URLParam(ctx.GetRequest(), "offset"))
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusBadRequest)
		return nil, err
	}
	request.Offset = offset
	limit, err := strconv.Atoi(chi.URLParam(ctx.GetRequest(), "limit"))
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusBadRequest)
		return nil, err
	}
	request.Limit = limit
	request.Query = ctx.GetRequest().URL.Query().Get("query")

	res, err := service.NewNewsService().GetNewsList(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return res, nil
}

func HandleGetNews(ctx context.Context) (any, error) {
	var request requests.GetNewsRequest
	request.ID = chi.URLParam(ctx.GetRequest(), "id")

	res, err := service.NewNewsService().GetNews(request)
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

func HandleCreateNews(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role != model.UserAdministrator {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.CreateNewsRequest
	decoder := json.NewDecoder(ctx.GetRequest().Body)
	err := decoder.Decode(&request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	err = service.NewNewsService().CreateNews(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}

func HandleDeleteNews(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role != model.UserAdministrator {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.DeleteNewsRequest
	request.ID = chi.URLParam(ctx.GetRequest(), "id")

	err := service.NewNewsService().DeleteNews(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}

func HandleUpdateNews(ctx context.Context) (any, error) {
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

func RouteNewsPaths(
	core handlers.Core,
	router chi.Router,
) {
	router.MethodFunc("GET", "/news", handlers.MakeHandler(handlers.MakeJSONWrapper(HandleGetNewsList), core))
	router.MethodFunc("GET", "/news/{id}", handlers.MakeHandler(handlers.MakeJSONWrapper(HandleGetNews), core))
	router.MethodFunc("PUT", "/news", handlers.MakeHandler(HandleCreateNews, core))
	router.MethodFunc("DELETE", "/news/{id}", handlers.MakeHandler(HandleDeleteNews, core))
	router.MethodFunc("PATCH", "/news/{id}", handlers.MakeHandler(HandleUpdateNews, core))
}
