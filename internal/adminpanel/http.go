package adminpanel

import (
	"netradio/internal/repository"
	"netradio/pkg/handlers"

	"github.com/go-chi/chi/v5"
)

///go:embed static
//var content embed.FS

func RoutePaths(
	core handlers.Core,
	router chi.Router,
	newsService repository.NewsDB,
	musicService repository.MusicDB,
) {

	addHandler(core, router, "POST", "/news/add", newCreateNewsHandler(newsService))
	addHandler(core, router, "POST", "/news/change", newChangeNewsHandler(newsService))
	addHandler(core, router, "POST", "/podcasts/add", newCreatePodcastHandler(musicService))
	addHandler(core, router, "POST", "/podcasts/change", newChangePodcastHandler(musicService))

	//fsys, _ := fs.Sub(content, "static")
	//router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(fsys))))
}

func addHandler(core handlers.Core, router chi.Router, method, pattern string, handler handlers.Handler) {
	//router.Method(method, pattern, handlers.NewAuthWrapper(handlers.NewResponseWrapper(handler), core))
}
