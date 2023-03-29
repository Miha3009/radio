package requests

import "time"

type GetChannelsRequest struct {
	Offset int
	Limit  int
	Query  string
	Status string
}

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

type UploadLogoRequest struct {
	ID   string
	Logo string
}

type GetCurrentTrackRequest struct {
	ID     string
	UserID *string
}

type GetScheduleRequest struct {
	ChannelID string
	Past      int
	Next      int
}

type GetScheduleRangeRequest struct {
	ChannelID string
	From      time.Time
	To        time.Time
}

type StartStreamRequest struct {
	ID string
}

type StopStreamRequest struct {
	ID string
}
