package model

import "time"

type ListenerStat struct {
	Count int       `json:"count"`
	Time  time.Time `json:"time"`
}
