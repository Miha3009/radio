package requests

type ListNewsRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
