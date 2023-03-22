package requests

type GetTrackRequest struct {
	ID     string
	UserID *string
}

type GetTrackListRequest struct {
	Offset int
	Limit  int
	Query  string
}

type CreateTrackRequest struct {
	Title       string `json:"title"`
	Perfomancer string `json:"description"`
	Year        int    `json:"year"`
}

type UpdateTrackRequest struct {
	ID          string
	Title       *string `json:"title,omitempty"`
	Perfomancer *string `json:"description,omitempty"`
	Year        *int    `json:"year,omitempty"`
}

type DeleteTrackRequest struct {
	ID string
}

type LikeTrackRequest struct {
	ID     string
	UserID int
	Like   bool `json:"like"`
}

type CommentTrackRequest struct {
	ID     string
	UserID int
	Text   string  `json:"text"`
	Parent *string `json:"parent,omitempty"`
}

type UploadTrackRequest struct {
	ID    string
	Audio string
}
