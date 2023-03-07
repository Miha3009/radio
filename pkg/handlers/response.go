package handlers

import "net/http"

type Response struct {
	Content    []byte
	StatusCode int
	Headers    http.Header
}

func (r *Response) GetStatusCodeOrDefault(defaultStatus int) int {
	if r == nil {
		return defaultStatus
	}
	if r.StatusCode < 100 {
		return defaultStatus
	}
	return r.StatusCode
}

func (r *Response) GetContent() []byte {
	if r == nil {
		return nil
	}
	return r.Content
}
