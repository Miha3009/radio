package handlers

import (
	"netradio/internal/repository"
	"netradio/pkg/log"
)

type Core struct {
	userService repository.UserDB
	log         log.Logger
}

func NewCore(log log.Logger, userService repository.UserDB) Core {
	return Core{
		log:         log,
		userService: userService,
	}
}
