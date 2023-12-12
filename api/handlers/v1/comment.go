package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golanguzb70/simple-post-app/models"
	"github.com/google/uuid"
)

// @Router		/comment [POST]
// @Summary		Create comment
// @Tags        Comment
// @Description	Here comment can be created.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       post   body       models.CommentCreateReq true "post info"
// @Success		200 	{object}  models.CommentResponse
// @Failure     default {object}  models.StandardResponse
func (h *handlerV1) CommentCreate(ctx *gin.Context) {
	claim, err := GetClaims(*h, ctx)
	if h.HandleResponse(ctx, err, http.StatusUnauthorized, UnAuthorized, "invalid authorization", nil) {
		return
	}

	body := &models.CommentCreateReq{}
	err = ctx.ShouldBindJSON(&body)
	if h.HandleResponse(ctx, err, http.StatusBadRequest, BadRequest, "invalid body", nil) {
		return
	}
	body.Id = uuid.New().String()
	body.UserId = claim.Sub

	res, err := h.storage.Postgres().CommentCreate(context.Background(), body)
	if h.HandleDatabaseLevelWithMessage(ctx, err, "CommentCreate: h.storage.Postgres().CommentCreate()") {
		return
	}

	h.HandleResponse(ctx, nil, http.StatusOK, Success, "", res)
}

// @Router		/comment/{id} [GET]
// @Summary		Get comment by key
// @Tags        Comment
// @Description	Here comment can be got.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       id       path     int true "id"
// @Success		200 	{object}  models.CommentResponse
// @Failure     default {object}  models.StandardResponse
func (h *handlerV1) CommentGet(ctx *gin.Context) {
	res, err := h.storage.Postgres().CommentGet(context.Background(), &models.CommentGetReq{
		Id: ctx.Param("id"),
	})

	if h.HandleDatabaseLevelWithMessage(ctx, err, "CommentGet: h.storage.Postgres().CommentGet()") {
		return
	}

	h.HandleResponse(ctx, nil, http.StatusOK, Success, "", res)
}

// @Router		/comment/list [GET]
// @Summary		Get comments list
// @Tags        Comment
// @Description	Here all comments can be got.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       filters query models.CommentFindReq true "filters"
// @Success		200 	{object}  models.CommentFindResponse
// @Failure     default {object}  models.StandardResponse
func (h *handlerV1) CommentFind(ctx *gin.Context) {
	var (
		dbReq = &models.CommentFindReq{}
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

	dbReq.PostId = ctx.Query("post_id")
	dbReq.OrderByCreatedAt, _ = strconv.ParseUint(ctx.Query("order_by_created_at"), 10, 8)

	res, err := h.storage.Postgres().CommentFind(context.Background(), dbReq)
	if h.HandleDatabaseLevelWithMessage(ctx, err, "CommentFind: h.storage.Postgres().CommentFind()") {
		return
	}

	h.HandleResponse(ctx, nil, http.StatusOK, Success, "", res)
}

// @Router		/comment [PUT]
// @Summary		Update comment
// @Tags        Comment
// @Description	Here comment can be updated.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       post   body       models.CommentUpdateReq true "post info"
// @Success		200 	{object}  models.CommentResponse
// @Failure     default {object}  models.StandardResponse
func (h *handlerV1) CommentUpdate(ctx *gin.Context) {
	claim, err := GetClaims(*h, ctx)
	if h.HandleResponse(ctx, err, http.StatusUnauthorized, UnAuthorized, "invalid authorization", nil) {
		return
	}

	body := &models.CommentUpdateReq{}
	err = ctx.ShouldBindJSON(&body)
	if h.HandleResponse(ctx, err, http.StatusBadRequest, BadRequest, "invalid body", nil) {
		return
	}

	body.UserId = claim.Sub
	res, err := h.storage.Postgres().CommentUpdate(context.Background(), body)
	if h.HandleDatabaseLevelWithMessage(ctx, err, "CommentUpdate: h.storage.Postgres().CommentUpdate()") {
		return
	}

	h.HandleResponse(ctx, nil, http.StatusOK, Success, "", res)
}

// @Router		/comment/{id} [DELETE]
// @Summary		Delete comment
// @Tags        Comment
// @Description	Here comment can be deleted.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       id       path     int true "id"
// @Success		200 	{object}  models.StandardResponse
// @Failure     default {object}  models.StandardResponse
func (h *handlerV1) CommentDelete(ctx *gin.Context) {
	claim, err := GetClaims(*h, ctx)
	if h.HandleResponse(ctx, err, http.StatusUnauthorized, UnAuthorized, "invalid authorization", nil) {
		return
	}

	err = h.storage.Postgres().CommentDelete(context.Background(), &models.CommentDeleteReq{Id: ctx.Param("id"), UserId: claim.Sub})
	if h.HandleDatabaseLevelWithMessage(ctx, err, "CommentDelete: h.storage.Postgres().CommentDelete()") {
		return
	}

	h.HandleResponse(ctx, nil, http.StatusOK, Success, "Successfully deleted", nil)
}
