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
	Title        string `json:"title"`
	Performancer string `json:"performancer"`
	Year         int    `json:"year"`
}

type UpdateTrackRequest struct {
	ID           string
	Title        *string `json:"title,omitempty"`
	Performancer *string `json:"performancer,omitempty"`
	Year         *int    `json:"year,omitempty"`
}

type DeleteTrackRequest struct {
	ID string
}

type LikeTrackRequest struct {
	ID     string
	UserID int
	Like   bool `json:"like"`
}

type GetTrackCommentsRequest struct {
	ID string
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
