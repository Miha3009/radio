package model

type ChannelInfo struct {
	ID          int           `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Logo        string        `json:"logo"`
	Status      ChannelStatus `json:"status"`
}

type ChannelStatus int

const (
	StoppedChannel ChannelStatus = 0
	ActiveChannel                = 1
)

type ChannelShortInfo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Logo  string `json:"logo"`
}
