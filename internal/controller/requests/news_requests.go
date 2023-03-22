package requests

type GetNewsListRequest struct {
	Offset int
	Limit  int
	Query  string
}

type GetNewsRequest struct {
	ID string
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
