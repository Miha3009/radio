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
	addStreamingHandler(core, router, "GET", "/podcast/{podcastID}/start", newPodcastGetter(musicService))
}

func addStreamingHandler(
	core handlers.Core,
	router chi.Router,
	method, pattern string,
	handler handlers.HandlerWritable,
) {
	//router.Method(method, pattern, handlers.NewAuthWrapper(handler, core))
}
