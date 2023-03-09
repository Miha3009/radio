package podcast

import (
	"netradio/internal/repository"
	"netradio/pkg/handlers"

	"github.com/go-chi/chi/v5"
)

func RoutePaths(
	core handlers.Core,
	router chi.Router,
	musicService repository.MusicDB,
) {
	//addJSONHandler(core, router, "GET", "/podcast/all", newGetterHandler(musicService))
	router.Method("POST", "/podcast/{podcastID}/start", handlers.MakeHandler(PodacstGetter, core))
}
