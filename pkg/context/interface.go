package context

import (
	"net/http"
	"netradio/internal/model"
	"netradio/pkg/log"
)

type Context interface {
	GetUser() model.User
	GetLogger() log.Logger
	GetRequest() *http.Request
	GetResponseWriter() http.ResponseWriter
}
