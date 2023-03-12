package responses

type GetTrackResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Perfomancer string `json:"perfomancer"`
	Year        int    `json:"year"`
	Liked       bool   `json:"liked"`
	Found       bool   `json:"-"`
}

type UpdateTrackResponse struct {
	Found bool `json:"-"`
}
