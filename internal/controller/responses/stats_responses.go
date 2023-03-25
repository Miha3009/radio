package responses

import "netradio/internal/model"

type GetListenersStatResponse struct {
	Stats []model.ListenerStat `json:"stats"`
}
