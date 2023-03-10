package responses

import "netradio/internal/model"

type GetChannelsResponse struct {
	Channels []model.ChannelShortInfo `json:"channels"`
}

type GetChannelResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Found       bool   `json:"-"`
}

type UpdateChannelResponse struct {
	Found bool `json:"-"`
}
