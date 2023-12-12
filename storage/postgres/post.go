package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/golanguzb70/simple-post-app/models"
)

func (r *postgresRepo) PostCreate(ctx context.Context, req *models.PostCreateReq) (*models.PostResponse, error) {
	res := &models.PostResponse{}
	query := r.Db.Builder.Insert("posts").Columns(
		"id, slug, title, content, user_id",
	).Values(req.Id, req.Slug, req.Title, req.Content, req.UserId).Suffix(
		"RETURNING id, slug, title, content, user_id, created_at, updated_at")

	err := query.RunWith(r.Db.Db).Scan(
		&res.Id, &res.Slug, &res.Title,
		&res.Content, &res.UserId,
		&CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "PostCreate: query.RunWith(r.Db.Db).Scan()")
	}

	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) PostGet(ctx context.Context, req *models.PostGetReq) (*models.PostResponse, error) {
	query := r.Db.Builder.Select("id, slug, title, content, user_id, created_at, updated_at").
		From("posts")

	if req.Id != "" {
		query = query.Where(squirrel.Eq{"id": req.Id})
	} else if req.Slug != "" {
		query = query.Where(squirrel.Eq{"slug": req.Slug})
	} else {
		return &models.PostResponse{}, fmt.Errorf("at least one filter should be exists")
	}

	res := &models.PostResponse{}
	err := query.RunWith(r.Db.Db).QueryRow().Scan(
		&res.Id, &res.Slug, &res.Title,
		&res.Content, &res.UserId,
		&CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "PostGet:query.RunWith(r.Db.Db).QueryRow()")
	}

	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) PostFind(ctx context.Context, req *models.PostFindReq) (*models.PostFindResponse, error) {
	var (
		res            = &models.PostFindResponse{}
		whereCondition = squirrel.And{}
		orderBy        = []string{}
	)

	if strings.TrimSpace(req.Search) != "" {
		whereCondition = append(whereCondition, squirrel.Or{squirrel.ILike{"slug": "%" + req.Search + "%"}, squirrel.ILike{"title": "%" + req.Search + "%"}})
	}

	if req.OrderByCreatedAt != 0 {
		if req.OrderByCreatedAt > 0 {
			orderBy = append(orderBy, "created_at DESC")
		} else {
			orderBy = append(orderBy, "created_at ASC")
		}
	}

	countQuery := r.Db.Builder.Select("count(1) as count").From("posts").Where("deleted_at is null").Where(whereCondition)
	err := countQuery.RunWith(r.Db.Db).QueryRow().Scan(&res.Count)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "PostFind: countQuery.RunWith(r.Db.Db).QueryRow().Scan()")
	}

	query := r.Db.Builder.Select("id, slug, title, content, user_id, created_at, updated_at").
		From("posts").Where("deleted_at is null").Where(whereCondition)

	if len(orderBy) > 0 {
		query = query.OrderBy(strings.Join(orderBy, ", "))
	}

	query = query.Limit(uint64(req.Limit)).Offset(uint64((req.Page - 1) * req.Limit))

	rows, err := query.RunWith(r.Db.Db).Query()
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "PostFind: query.RunWith(r.Db.Db).Query()")
	}
	defer rows.Close()

	for rows.Next() {
		temp := &models.PostResponse{}
		err := rows.Scan(
			&temp.Id, &temp.Slug,
			&temp.Title, &temp.Content, &temp.UserId,
			&CreatedAt, &UpdatedAt,
		)
		if err != nil {
			return res, HandleDatabaseError(err, r.Log, "PostFind: rows.Scan()")
		}

		temp.CreatedAt = CreatedAt.Format(time.RFC1123)
		temp.UpdatedAt = UpdatedAt.Format(time.RFC1123)
		res.Posts = append(res.Posts, temp)
	}

	return res, nil
}

func (r *postgresRepo) PostUpdate(ctx context.Context, req *models.PostUpdateReq) (*models.PostResponse, error) {
	var (
		mp             = make(map[string]interface{})
		whereCondition = squirrel.And{squirrel.Eq{"id": req.Id}}
	)

	mp["content"] = req.Content
	mp["title"] = req.Title
	mp["updated_at"] = time.Now()

	if req.UserId != "" {
		whereCondition = append(whereCondition, squirrel.Eq{"user_id": req.UserId})
	}

	query := r.Db.Builder.Update("posts").SetMap(mp).
		Where(whereCondition).
		Suffix("RETURNING id, slug, title, content, user_id, created_at, updated_at")

	res := &models.PostResponse{}
	err := query.RunWith(r.Db.Db).QueryRow().Scan(
		&res.Id, &res.Slug,
		&res.Title, &res.Content,
		&res.UserId,
		&CreatedAt, &UpdatedAt,
	)

	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "PostUpdate: query.RunWith(r.Db.Db).QueryRow().Scan()")
	}

	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) PostDelete(ctx context.Context, req *models.PostDeleteReq) error {
	whereCondition := squirrel.And{squirrel.Eq{"id": req.Id}}

	if req.UserId != "" {
		whereCondition = append(whereCondition, squirrel.Eq{"user_id": req.UserId})
	}

	query := r.Db.Builder.Delete("posts").Where(whereCondition)

	_, err := query.RunWith(r.Db.Db).Exec()
	return HandleDatabaseError(err, r.Log, "PostDelete: query.RunWith(r.Db.Db).Exec()")
}
