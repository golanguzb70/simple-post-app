package postgres

import (
	"context"

	"github.com/golanguzb70/simple-post-app/models"
)

type PostgresI interface {
	// common
	UpdateSingleField(ctx context.Context, req *models.UpdateSingleFieldReq) error
	CheckIfExists(ctx context.Context, req *models.CheckIfExistsReq) (*models.CheckIfExistsRes, error)

	// User
	UserCreate(ctx context.Context, req *models.UserCreateReq) (*models.UserResponse, error)
	UserGet(ctx context.Context, req *models.UserGetReq) (*models.UserResponse, error)
	UserFind(ctx context.Context, req *models.UserFindReq) (*models.UserFindResponse, error)
	UserUpdate(ctx context.Context, req *models.UserUpdateReq) (*models.UserResponse, error)
	UserDelete(ctx context.Context, req *models.UserDeleteReq) error

	// Post 
	PostCreate(ctx context.Context, req *models.PostCreateReq) (*models.PostResponse, error)
	PostGet(ctx context.Context, req *models.PostGetReq) (*models.PostResponse, error)
	PostFind(ctx context.Context, req *models.PostFindReq) (*models.PostFindResponse, error)
	PostUpdate(ctx context.Context, req *models.PostUpdateReq) (*models.PostResponse, error)
	PostDelete(ctx context.Context, req *models.PostDeleteReq) error
	

   // Comment
	CommentCreate(ctx context.Context, req *models.CommentCreateReq) (*models.CommentResponse, error)
	CommentGet(ctx context.Context, req *models.CommentGetReq) (*models.CommentResponse, error)
	CommentFind(ctx context.Context, req *models.CommentFindReq) (*models.CommentFindResponse, error)
	CommentUpdate(ctx context.Context, req *models.CommentUpdateReq) (*models.CommentResponse, error)
	CommentDelete(ctx context.Context, req *models.CommentDeleteReq) error
	// Don't delete this line, it is used to modify the file automatically
}
