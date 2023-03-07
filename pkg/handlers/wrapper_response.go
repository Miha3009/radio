package handlers

import (
	"net/http"
	"netradio/pkg/context"
)

func NewResponseWrapper(handler Handler) *responseWrapper {
	return &responseWrapper{
		original: handler,
	}
}

type responseWrapper struct {
	original Handler
}

func (w *responseWrapper) ServeHTTP(ctx context.Context, responseWriter http.ResponseWriter) error {
	response, err := w.original.ServeHTTP(ctx)

	for k, v := range response.Headers {
		if len(v) == 0 {
			continue
		}
		responseWriter.Header().Add(k, v[0])
	}

	if err != nil {
		return err
	}

	statusCode := response.GetStatusCodeOrDefault(http.StatusOK)
	responseWriter.WriteHeader(statusCode)

	if response.Content == nil {
		return nil
	}
	_, err = responseWriter.Write(response.GetContent())
	if err != nil {
		return err
	}
	return nil
}
