package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golanguzb70/simple-post-app/models"
	"github.com/google/uuid"
)

// @Router		/post [POST]
// @Summary		Create post
// @Tags        Post
// @Description	Here post can be created.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       post   body       models.PostCreateReq true "post info"
// @Success		200 	{object}  models.PostResponse
// @Failure     default {object}  models.StandardResponse
func (h *handlerV1) PostCreate(ctx *gin.Context) {
	claim, err := GetClaims(*h, ctx)
	if h.HandleResponse(ctx, err, http.StatusUnauthorized, UnAuthorized, "invalid authorization", nil) {
		return
	}

	body := &models.PostCreateReq{}
	err = ctx.ShouldBindJSON(&body)
	if h.HandleResponse(ctx, err, http.StatusBadRequest, BadRequest, "invalid body", nil) {
		return
	}

	body.UserId = claim.Sub
	body.Id = uuid.New().String()

	res, err := h.storage.Postgres().PostCreate(context.Background(), body)
	if h.HandleDatabaseLevelWithMessage(ctx, err, "PostCreate: h.storage.Postgres().PostCreate()") {
		return
	}

	h.HandleResponse(ctx, nil, http.StatusCreated, Success, "", res)
}

// @Router		/post/{id} [GET]
// @Summary		Get post by key
// @Tags        Post
// @Description	Here post can be got.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       id       path     string true "id"
// @Success		200 	{object}  models.PostResponse
// @Failure     default {object}  models.StandardResponse
func (h *handlerV1) PostGet(ctx *gin.Context) {
	res, err := h.storage.Postgres().PostGet(context.Background(), &models.PostGetReq{
		Id: ctx.Param("id"),
	})

	if h.HandleDatabaseLevelWithMessage(ctx, err, "PostGet: h.storage.Postgres().PostGet()") {
		return
	}

	h.HandleResponse(ctx, nil, http.StatusOK, Success, "", res)
}

// @Router		/post/list [GET]
// @Summary		Get posts list
// @Tags        Post
// @Description	Here all posts can be got.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       filters query models.PostFindReq true "filters"
// @Success		200 	{object}  models.PostFindResponse
// @Failure     default {object}  models.StandardResponse
func (h *handlerV1) PostFind(ctx *gin.Context) {
	var (
		dbReq = &models.PostFindReq{}
		err   error
	)

	dbReq.Page, err = ParsePageQueryParam(ctx)
	if h.HandleResponse(ctx, err, http.StatusBadRequest, BadRequest, "invalid page param", nil) {
		return
	}

	dbReq.Limit, err = ParseLimitQueryParam(ctx)
	if h.HandleResponse(ctx, err, http.StatusBadRequest, BadRequest, "invalid limit param", nil) {
		return
	}

	dbReq.Search = ctx.Query("search")
	dbReq.OrderByCreatedAt, _ = strconv.ParseUint(ctx.Query("order_by_created_at"), 10, 8)

	res, err := h.storage.Postgres().PostFind(context.Background(), dbReq)
	if h.HandleDatabaseLevelWithMessage(ctx, err, "PostFind: h.storage.Postgres().PostFind()") {
		return
	}

	h.HandleResponse(ctx, nil, http.StatusOK, Success, "", res)
}

// @Router		/post [PUT]
// @Summary		Update post
// @Tags        Post
// @Description	Here post can be updated.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       post   body       models.PostUpdateReq true "post info"
// @Success		200 	{object}  models.PostResponse
// @Failure     default {object}  models.StandardResponse
func (h *handlerV1) PostUpdate(ctx *gin.Context) {
	claim, err := GetClaims(*h, ctx)
	if h.HandleResponse(ctx, err, http.StatusUnauthorized, UnAuthorized, "invalid authorization", nil) {
		return
	}

	body := &models.PostUpdateReq{}
	err = ctx.ShouldBindJSON(&body)
	if h.HandleResponse(ctx, err, http.StatusBadRequest, BadRequest, "invalid body", nil) {
		return
	}

	body.UserId = claim.Sub

	res, err := h.storage.Postgres().PostUpdate(context.Background(), body)
	if h.HandleDatabaseLevelWithMessage(ctx, err, "PostUpdate: h.storage.Postgres().PostUpdate()") {
		return
	}

	h.HandleResponse(ctx, nil, http.StatusOK, Success, "", res)
}

// @Router		/post/{id} [DELETE]
// @Summary		Delete post
// @Tags        Post
// @Description	Here post can be deleted.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       id       path     int true "id"
// @Success		200 	{object}  models.StandardResponse
// @Failure     default {object}  models.StandardResponse
func (h *handlerV1) PostDelete(ctx *gin.Context) {
	claim, err := GetClaims(*h, ctx)
	if h.HandleResponse(ctx, err, http.StatusUnauthorized, UnAuthorized, "invalid authorization", nil) {
		return
	}

	err = h.storage.Postgres().PostDelete(context.Background(), &models.PostDeleteReq{Id: ctx.Param("id"), UserId: claim.Sub})
	if h.HandleDatabaseLevelWithMessage(ctx, err, "PostDelete: h.storage.Postgres().PostDelete()") {
		return
	}

	h.HandleResponse(ctx, nil, http.StatusOK, Success, "Successfully deleted", nil)
}
