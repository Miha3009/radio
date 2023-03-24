package model

import "time"

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

type ScheduleTrack struct {
	ID        string     `json:"id"`
	TrackId   *string    `json:"trackid"`
	ChannelId *string    `json:"channelid"`
	StartDate *time.Time `json:"startdate"`
	EndDate   *time.Time `json:"enddate"`
}

type ScheduleTrackFull struct {
	ID        string     `json:"id"`
	Track     Track      `json:"trackid"`
	ChannelId string     `json:"channelid"`
	StartDate *time.Time `json:"startdate"`
	EndDate   *time.Time `json:"enddate"`
}
