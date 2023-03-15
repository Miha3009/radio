package podcast

import (
	"netradio/pkg/handlers"

	"github.com/go-chi/chi/v5"
)

func RoutePaths(
	core handlers.Core,
	router chi.Router,
) {
	//addJSONHandler(core, router, "GET", "/podcast/all", newGetterHandler(musicService))
	router.Method("POST", "/podcast/{podcastID}/start", handlers.MakeHandler(PodacstGetter, core))
}
