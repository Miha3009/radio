package handlers

import (
	"encoding/json"
	"net/http"
	"netradio/pkg/context"
)

type Handler interface {
	ServeHTTP(context.Context) (Response, error)
}

type HandlerWritable interface {
	ServeHTTP(context.Context, http.ResponseWriter) error
}

type HandlerJSON interface {
	ServeJSON(context context.Context) (json.RawMessage, error)
}

type HandlerFunc func(context.Context) (any, error)

func MakeHandler(handler HandlerFunc, core Core) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, r *http.Request) {
		ctx := context.NewContext(r, responseWriter, core.userService, core.log)
		_, err := handler(ctx)
		if err != nil {
			core.log.Error(err)
		}
	}
}
