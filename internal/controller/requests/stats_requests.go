package requests

import "time"

type GetListenersStatRequest struct {
	ChannelID string
	From      time.Time
	To        time.Time
}
