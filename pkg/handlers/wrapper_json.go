package handlers

import (
	"encoding/json"
	"netradio/pkg/context"
)

func MakeJSONWrapper(handler HandlerFunc) HandlerFunc {
	return func(context context.Context) (any, error) {
		object, err := handler(context)
		if err != nil {
			return nil, err
		}

		res, err := json.Marshal(object)
		if err != nil {
			return nil, err
		}

		responseWriter := context.GetResponseWriter()
		_, err = responseWriter.Write(res)
		if err != nil {
			return nil, err
		}

		responseWriter.Header().Set("Content-Type", "application/json")

		return res, nil
	}
}
