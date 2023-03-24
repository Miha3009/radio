package requests

import "time"

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

type AddTrackRequest struct {
	ID        string
	TrackID   string    `json:"trackid"`
	StartDate time.Time `json:"startdate"`
}

type GetCurrentTrackRequest struct {
	ID     string
	UserID *string
}

type GetScheduleRequest struct {
	ID   string
	Past int
	Next int
}
