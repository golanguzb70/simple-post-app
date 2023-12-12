package models

type CommentCreateReq struct {
	Id      string `json:"id"`
	UserId  string `json:"user_id"`
	PostId  string `json:"post_id"`
	Comment string `json:"comment"`
}

type CommentUpdateReq struct {
	Id      string `json:"id"`
	UserId  string `json:"user_id"`
	Comment string `json:"comment"`
}

type CommentGetReq struct {
	Id string `json:"id"`
}

type CommentFindReq struct {
	Page             int    `json:"page"`
	Limit            int    `json:"limit"`
	OrderByCreatedAt uint64 `json:"order_by_created_at"`
	PostId           string `json:"post_id"`
}

type CommentDeleteReq struct {
	Id     string `json:"id"`
	UserId string `json:"user_id"`
}

type CommentFindResponse struct {
	Comments []*CommentResponse `json:"comments"`
	Count    int                `json:"count"`
}

type CommentResponse struct {
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	PostId    string `json:"post_id"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
