package requests

type GetNewsListRequest struct {
	Offset int
	Limit  int
	Query  string
}

type GetNewsRequest struct {
	ID     string
	UserID *string
}

type CreateNewsRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateNewsRequest struct {
	ID      string
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

type DeleteNewsRequest struct {
	ID string
}

type UploadImageRequest struct {
	ID    string
	Image string
}

type LikeNewsRequest struct {
	ID     string
	UserID int
	Like   bool `json:"like"`
}

type GetNewsCommentsRequest struct {
	ID string
}

type CommentNewsRequest struct {
	ID     string
	UserID int
	Text   string  `json:"text"`
	Parent *string `json:"parent,omitempty"`
}
