package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/golanguzb70/simple-post-app/models"
)

func (r *postgresRepo) CommentCreate(ctx context.Context, req *models.CommentCreateReq) (*models.CommentResponse, error) {
	res := &models.CommentResponse{}
	query := r.Db.Builder.Insert("comments").Columns(
		"id, user_id, post_id, comment",
	).Values(req.Id, req.UserId, req.PostId, req.Comment).Suffix(
		"RETURNING id, user_id, post_id, comment, created_at, updated_at")

	err := query.RunWith(r.Db.Db).Scan(
		&res.Id, &res.UserId, &res.PostId, &res.Comment,
		&CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "CommentCreate: query.RunWith(r.Db.Db).Scan()")
	}

	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) CommentGet(ctx context.Context, req *models.CommentGetReq) (*models.CommentResponse, error) {
	query := r.Db.Builder.Select("id, user_id, post_id, comment, created_at, updated_at").
		From("comments")

	if req.Id != "" {
		query = query.Where(squirrel.Eq{"id": req.Id})
	} else {
		return &models.CommentResponse{}, fmt.Errorf("at least one filter should be exists")
	}

	res := &models.CommentResponse{}
	err := query.RunWith(r.Db.Db).QueryRow().Scan(
		&res.Id, &res.UserId, &res.PostId, &res.Comment,
		&CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "CommentGet:query.RunWith(r.Db.Db).QueryRow()")
	}

	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) CommentFind(ctx context.Context, req *models.CommentFindReq) (*models.CommentFindResponse, error) {
	var (
		res            = &models.CommentFindResponse{}
		whereCondition = squirrel.And{}
		orderBy        = []string{}
	)

	if strings.TrimSpace(req.PostId) != "" {
		whereCondition = append(whereCondition, squirrel.Eq{"post_id": req.PostId})
	}

	if req.OrderByCreatedAt != 0 {
		if req.OrderByCreatedAt > 0 {
			orderBy = append(orderBy, "created_at DESC")
		} else {
			orderBy = append(orderBy, "created_at ASC")
		}
	}

	countQuery := r.Db.Builder.Select("count(1) as count").From("comments").Where("deleted_at is null").Where(whereCondition)
	err := countQuery.RunWith(r.Db.Db).QueryRow().Scan(&res.Count)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "CommentFind: countQuery.RunWith(r.Db.Db).QueryRow().Scan()")
	}

	query := r.Db.Builder.Select("id, user_id, post_id, comment, created_at, updated_at").
		From("comments").Where("deleted_at is null").Where(whereCondition)

	if len(orderBy) > 0 {
		query = query.OrderBy(strings.Join(orderBy, ", "))
	}

	query = query.Limit(uint64(req.Limit)).Offset(uint64((req.Page - 1) * req.Limit))

	rows, err := query.RunWith(r.Db.Db).Query()
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "CommentFind: query.RunWith(r.Db.Db).Query()")
	}
	defer rows.Close()

	for rows.Next() {
		temp := &models.CommentResponse{}
		err := rows.Scan(
			&temp.Id, &temp.UserId, &temp.PostId, &temp.Comment,
			&CreatedAt, &UpdatedAt,
		)
		if err != nil {
			return res, HandleDatabaseError(err, r.Log, "CommentFind: rows.Scan()")
		}

		temp.CreatedAt = CreatedAt.Format(time.RFC1123)
		temp.UpdatedAt = UpdatedAt.Format(time.RFC1123)
		res.Comments = append(res.Comments, temp)
	}

	return res, nil
}

func (r *postgresRepo) CommentUpdate(ctx context.Context, req *models.CommentUpdateReq) (*models.CommentResponse, error) {
	var (
		mp             = make(map[string]interface{})
		whereCondition = squirrel.And{squirrel.Eq{"id": req.Id}}
	)

	mp["comment"] = req.Comment
	mp["updated_at"] = time.Now()

	if req.UserId != "" {
		whereCondition = append(whereCondition, squirrel.Eq{"user_id": req.UserId})
	}

	query := r.Db.Builder.Update("comments").SetMap(mp).
		Where(whereCondition).
		Suffix("RETURNING id, user_id, post_id, comment, created_at, updated_at")

	res := &models.CommentResponse{}
	err := query.RunWith(r.Db.Db).QueryRow().Scan(
		&res.Id, &res.UserId, &res.PostId, &res.Comment,
		&CreatedAt, &UpdatedAt,
	)

	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "CommentUpdate: query.RunWith(r.Db.Db).QueryRow().Scan()")
	}
	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) CommentDelete(ctx context.Context, req *models.CommentDeleteReq) error {
	whereCondition := squirrel.And{squirrel.Eq{"id": req.Id}}

	if req.UserId != "" {
		whereCondition = append(whereCondition, squirrel.Eq{"user_id": req.UserId})
	}

	query := r.Db.Builder.Delete("comments").Where(whereCondition)

	_, err := query.RunWith(r.Db.Db).Exec()
	return HandleDatabaseError(err, r.Log, "CommentDelete: query.RunWith(r.Db.Db).Exec()")
}
