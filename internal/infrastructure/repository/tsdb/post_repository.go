package tsdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/minipkg/selection_condition"

	"blog/internal/pkg/apperror"

	"blog/internal/domain/post"
)

type PostRepository struct {
	*Repository
}

const ()

var _ post.WriteRepository = (*PostRepository)(nil)
var _ post.ReadRepository = (*PostRepository)(nil)

func NewPostRepository(repository *Repository) *PostRepository {
	return &PostRepository{
		Repository: repository,
	}
}

const (
	post_sql_Get          = "SELECT id, blog_id, sysname, keyword_ids, tag_ids, is_deleted, title, preview, content, created_at, updated_at, deleted_at FROM blog.post WHERE id = $1;"
	post_sql_GetBySysname = "SELECT id, blog_id, sysname, keyword_ids, tag_ids, is_deleted, title, preview, content, created_at, updated_at, deleted_at FROM blog.post WHERE sysname = $1;"
	post_sql_MGet         = "SELECT id, blog_id, sysname, keyword_ids, tag_ids, is_deleted, title, preview, created_at, updated_at, deleted_at FROM blog.post WHERE id = any($1);"
	post_sql_Filter       = "SELECT id, blog_id, sysname, keyword_ids, tag_ids, is_deleted, title, preview, created_at, updated_at, deleted_at FROM blog.post;"
	post_sql_Create       = "INSERT INTO blog.post(blog_id, sysname, keyword_ids, tag_ids, title, preview, content, content_tsvector, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, to_tsvector('russian', $7), Now()) RETURNING id;"
	post_sql_Update       = "UPDATE blog.post SET blog_id = $2, sysname = $3, keyword_ids = $4, tag_ids = $5, title = $6, preview = $7, content = $8, content_tsvector= to_tsvector('russian', $8), updated_at = Now() WHERE id = $1;"
	post_sql_Delete       = "DELETE FROM blog.post WHERE id = $1;"
)

func (r *PostRepository) Get(ctx context.Context, ID uint) (*post.Post, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "PostRepository.Get"
	start := time.Now().UTC()

	entity := &post.Post{}
	if err := r.db.QueryRow(ctx, post_sql_Get, ID).Scan(&entity.ID, &entity.BlogID, &entity.Sysname, &entity.KeywordIDs, &entity.TagIDs, &entity.IsDeleted, &entity.Title, &entity.Preview, &entity.Content, &entity.CreatedAt, &entity.UpdatedAt, &entity.DeletedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, post_sql_Get, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return entity, nil
}

func (r *PostRepository) GetBySysname(ctx context.Context, sysname string) (*post.Post, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "PostRepository.GetBySysname"
	start := time.Now().UTC()

	entity := &post.Post{}
	if err := r.db.QueryRow(ctx, post_sql_GetBySysname, sysname).Scan(&entity.ID, &entity.BlogID, &entity.Sysname, &entity.KeywordIDs, &entity.TagIDs, &entity.IsDeleted, &entity.Title, &entity.Preview, &entity.Content, &entity.CreatedAt, &entity.UpdatedAt, &entity.DeletedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, post_sql_GetBySysname, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return entity, nil
}

func (r *PostRepository) MGet(ctx context.Context, IDs *[]uint) (*[]post.PostPreview, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "PostRepository.MGet"

	var item post.PostPreview
	res := make([]post.PostPreview, 0, len(*IDs))

	start := time.Now().UTC()
	rows, err := r.db.Query(ctx, post_sql_MGet, *IDs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, post_sql_MGet, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&item.ID, &item.BlogID, &item.Sysname, &item.KeywordIDs, &item.TagIDs, &item.IsDeleted, &item.Title, &item.Preview, &item.CreatedAt, &item.UpdatedAt, &item.DeletedAt); err != nil {
			r.metrics.SqlMetrics.Inc(metricName, metricsFail)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
			return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, post_sql_MGet, err)
		}
		res = append(res, item)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)

	if len(res) == 0 {
		return nil, apperror.ErrNotFound
	}

	return &res, nil
}

func (r *PostRepository) Filter(ctx context.Context, condition *selection_condition.SelectionCondition) (*[]post.PostPreview, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "PostRepository.Filter"

	var item post.PostPreview
	res := make([]post.PostPreview, 0, defaultCapacityForResult)

	start := time.Now().UTC()
	rows, err := r.db.Query(ctx, post_sql_Filter)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, post_sql_Filter, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&item.ID, &item.BlogID, &item.Sysname, &item.KeywordIDs, &item.TagIDs, &item.IsDeleted, &item.Title, &item.Preview, &item.CreatedAt, &item.UpdatedAt, &item.DeletedAt); err != nil {
			r.metrics.SqlMetrics.Inc(metricName, metricsFail)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
			return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, post_sql_Filter, err)
		}
		res = append(res, item)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)

	if len(res) == 0 {
		return nil, apperror.ErrNotFound
	}

	return &res, nil
}

func (r *PostRepository) Create(ctx context.Context, entity *post.Post) (ID uint, err error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	const metricName = "PostRepository.Create"
	start := time.Now().UTC()

	if err := r.db.QueryRow(ctx, post_sql_Create, entity.BlogID, entity.Sysname, entity.KeywordIDs, entity.TagIDs, entity.IsDeleted, entity.Title, entity.Preview, entity.Content, entity.CreatedAt, entity.UpdatedAt, entity.DeletedAt).Scan(&ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return 0, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return 0, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, post_sql_Create, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return ID, nil
}

func (r *PostRepository) Update(ctx context.Context, entity *post.Post) error {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "PostRepository.Update"
	start := time.Now().UTC()

	_, err := r.db.Exec(ctx, post_sql_Update, entity.ID, entity.BlogID, entity.Sysname, entity.KeywordIDs, entity.TagIDs, entity.IsDeleted, entity.Title, entity.Preview, entity.Content, entity.CreatedAt, entity.UpdatedAt, entity.DeletedAt)
	if err != nil {
		if strings.Contains(err.Error(), errMsg_duplicateKey) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return apperror.ErrBadRequest
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, post_sql_Update, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}

func (r *PostRepository) Delete(ctx context.Context, ID uint) error {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "PostRepository.Delete"
	start := time.Now().UTC()

	_, err := r.db.Exec(ctx, post_sql_Delete, ID)
	if err != nil {
		if strings.Contains(err.Error(), errMsg_duplicateKey) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return apperror.ErrBadRequest
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, post_sql_Delete, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}

func (r *PostRepository) TextSearch(ctx context.Context, searchString string, createdAtSortOrder *string) (*[]post.PostPreview, error) {
	return nil, nil
}
