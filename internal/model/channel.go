package model

type ChannelInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ChannelSchedule struct {
	TimeStart int       `json:"time_start"`
	Duration  int       `json:"duration"`
	MusicInfo MusicInfo `json:"music_info"`
}
