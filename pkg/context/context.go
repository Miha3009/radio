package context

import (
	"net/http"
	"netradio/internal/model"
	"netradio/internal/repository"
	"netradio/pkg/jwt"
	"netradio/pkg/log"
	"strconv"
)

func NewContext(request *http.Request, responseWriter http.ResponseWriter, userService repository.UserDB, logger log.Logger) *ContextImpl {
	return &ContextImpl{
		request:        request,
		responseWriter: responseWriter,
		userService:    userService,
		logger:         logger,
	}
}

type ContextImpl struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	userService    repository.UserDB
	logger         log.Logger
}

func (c *ContextImpl) GetUser() model.User {
	userID, err := jwt.GetUserIDFromHeader(c.request.Header)
	if err != nil {
		return makeDefaultUser()
	}
	user, err := repository.NewUserDB().GetUserById(strconv.Itoa(userID))
	if user == nil {
		return makeDefaultUser()
	}
	return *user
}

func (c *ContextImpl) GetLogger() log.Logger {
	return c.logger
}

func (c *ContextImpl) GetRequest() *http.Request {
	return c.request
}

func (c *ContextImpl) GetResponseWriter() http.ResponseWriter {
	return c.responseWriter
}

func makeDefaultUser() model.User {
	return model.User{Role: model.UserGuest}
}
