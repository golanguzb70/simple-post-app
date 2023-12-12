package models

type PostCreateReq struct {
	Id      string `json:"id"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserId  string `json:"user_id"`
}

type PostUpdateReq struct {
	Id      string `json:"id"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserId  string `json:"user_id"`
}

type PostGetReq struct {
	Id   string `json:"id"`
	Slug string `json:"slug"`
}

type PostFindReq struct {
	Page             int    `json:"page"`
	Limit            int    `json:"limit"`
	OrderByCreatedAt uint64 `json:"order_by_created_at"`
	Search           string `json:"search"`
}

type PostDeleteReq struct {
	Id     string `json:"id"`
	UserId string `json:"user_id"`
}

type PostFindResponse struct {
	Posts []*PostResponse `json:"posts"`
	Count int             `json:"count"`
}

type PostResponse struct {
	Id        string `json:"id"`
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	UserId    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
