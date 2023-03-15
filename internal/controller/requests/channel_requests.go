package requests

type GetChannelRequest struct {
	ID string
}

type CreateChannelRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateChannelRequest struct {
	ID          string
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
}

type DeleteChannelRequest struct {
	ID string
}

type StartChannelRequest struct {
	ID string
}

type StopChannelRequest struct {
	ID string
}

type ConnectChannelRequest struct {
	ID  string
	SDP string `json:"sdp"`
}
