package model

type ChannelInfo struct {
	ID          int           `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Status      ChannelStatus `json:"status"`
}

type ChannelStatus int

const (
	ActiveChannel  ChannelStatus = 0
	StoppedChannel               = 1
)

type ChannelShortInfo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}
