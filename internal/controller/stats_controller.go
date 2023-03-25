package controller

import (
	"net/http"
	"netradio/internal/controller/requests"
	"netradio/internal/model"
	"netradio/internal/service"
	"netradio/pkg/context"
	"netradio/pkg/handlers"
	"time"

	"github.com/go-chi/chi/v5"
)

func HandleGetListenersStat(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role != model.UserAdministrator {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.GetListenersStatRequest
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

	res, err := service.NewStatsService().GetListenersStat(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return res, nil
}

func RouteStatsPaths(
	core handlers.Core,
	router chi.Router,
) {
	router.MethodFunc("GET", "/stats/channel/{id}/listeners", handlers.MakeHandler(handlers.MakeJSONWrapper(HandleGetListenersStat), core))
}
