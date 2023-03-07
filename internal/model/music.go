package model

type MusicInfo struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	PhotoLink string        `json:"photo_link"`
	Category  MusicCategory `json:"category"`
}

type MusicCategory = string

const (
	Podcast     MusicCategory = "podcast"
	Music       MusicCategory = "music"
	NewsProgram MusicCategory = "news_program"
)
