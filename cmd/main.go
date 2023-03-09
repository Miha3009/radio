package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"netradio/pkg/handlers"
	"netradio/pkg/jwt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"netradio/internal/adminpanel"
	"netradio/internal/controller"
	"netradio/internal/podcast"
	"netradio/internal/repository"
	"netradio/pkg/config"
	"netradio/pkg/database"
	"netradio/pkg/email"
	"netradio/pkg/errors"
	"netradio/pkg/log"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	logger := log.NewLogger()

	cfg, err := config.NewConfigFromYAML(config.DefaultYAMLPath)
	if err != nil {
		logger.Fatal(err)
	}

	err = database.OpenConnection(cfg.Database)
	if err != nil {
		logger.Fatal(err)
	}
	defer database.CloseConnection()

	jwt.SetConfig(cfg.Jwt)
	email.SetConfig(cfg.Email)

	newsDB := repository.NewNewsDB()
	musicDB := repository.NewMusicDB()
	userDB := repository.NewUserDB()

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	core := handlers.NewCore(logger, userDB)
	adminpanel.RoutePaths(core, router, newsDB, musicDB)
	controller.RouteUserPaths(core, router)
	controller.RouteNewsPaths(core, router)
	podcast.RoutePaths(core, router, musicDB)

	server := &http.Server{}
	server.Addr = fmt.Sprintf(":%d", cfg.Port)
	server.Handler = router

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(errors.Wrap(err, "http server failure"))
			sigChan <- syscall.SIGINT
		}
	}()

	<-sigChan
	logger.Info("shutting down")

}
