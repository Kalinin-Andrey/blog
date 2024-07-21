package tsdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"blog/internal/pkg/apperror"

	"blog/internal/domain/blog"
)

type BlogRepository struct {
	*Repository
}

const ()

var _ blog.WriteRepository = (*BlogRepository)(nil)
var _ blog.ReadRepository = (*BlogRepository)(nil)

func NewBlogRepository(repository *Repository) *BlogRepository {
	return &BlogRepository{
		Repository: repository,
	}
}

const (
	blog_sql_Get    = "SELECT id, sysname, keyword_ids, tag_ids, name, description FROM blog.blog WHERE sysname = $1;"
	blog_sql_MGet   = "SELECT id, sysname, keyword_ids, tag_ids, name, description FROM blog.blog WHERE sysname = any($1);"
	blog_sql_GetAll = "SELECT id, sysname, keyword_ids, tag_ids, name, description FROM blog.blog;"
	blog_sql_Create = "INSERT INTO blog.blog(sysname, keyword_ids, tag_ids, name, description) VALUES ($1, $2, $3, $4, $5) RETURNING id;"
	blog_sql_Update = "UPDATE blog.blog SET sysname = $2, keyword_ids = $3, tag_ids = $4, name = $5, description = $3 WHERE id = $1;"
	blog_sql_Delete = "DELETE FROM blog.blog WHERE id = $1;"
)

func (r *BlogRepository) Get(ctx context.Context, sysname string) (*blog.Blog, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "BlogRepository.Get"
	start := time.Now().UTC()

	entity := &blog.Blog{}
	if err := r.db.QueryRow(ctx, blog_sql_Get, sysname).Scan(&entity.ID, &entity.Sysname, &entity.KeywordIDs, &entity.TagIDs, &entity.Name, &entity.Description); err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, blog_sql_Get, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return entity, nil
}

func (r *BlogRepository) MGet(ctx context.Context, sysnames *[]string) (*[]blog.Blog, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "BlogRepository.MGet"

	var item blog.Blog
	res := make([]blog.Blog, 0, len(*sysnames))

	start := time.Now().UTC()
	rows, err := r.db.Query(ctx, blog_sql_MGet, *sysnames)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, blog_sql_MGet, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&item.ID, &item.Sysname, &item.KeywordIDs, &item.TagIDs, &item.Name, &item.Description); err != nil {
			r.metrics.SqlMetrics.Inc(metricName, metricsFail)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
			return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, blog_sql_MGet, err)
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

func (r *BlogRepository) GetAll(ctx context.Context) (*[]blog.Blog, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "BlogRepository.GetAll"

	var item blog.Blog
	res := make([]blog.Blog, 0, defaultCapacityForResult)

	start := time.Now().UTC()
	rows, err := r.db.Query(ctx, blog_sql_GetAll)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, blog_sql_GetAll, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&item.ID, &item.Sysname, &item.KeywordIDs, &item.TagIDs, &item.Name, &item.Description); err != nil {
			r.metrics.SqlMetrics.Inc(metricName, metricsFail)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
			return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, blog_sql_GetAll, err)
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

func (r *BlogRepository) Create(ctx context.Context, entity *blog.Blog) (ID uint, err error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	const metricName = "BlogRepository.Create"
	start := time.Now().UTC()

	if err := r.db.QueryRow(ctx, blog_sql_Create, entity.Sysname, entity.KeywordIDs, entity.TagIDs, entity.Name, entity.Description).Scan(&ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return 0, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return 0, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, blog_sql_Create, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return ID, nil
}

func (r *BlogRepository) Update(ctx context.Context, entity *blog.Blog) error {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "BlogRepository.Update"
	start := time.Now().UTC()

	_, err := r.db.Exec(ctx, blog_sql_Update, entity.ID, entity.Sysname, entity.KeywordIDs, entity.TagIDs, entity.Name, entity.Description)
	if err != nil {
		if strings.Contains(err.Error(), errMsg_duplicateKey) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return apperror.ErrBadRequest
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, blog_sql_Update, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}

func (r *BlogRepository) Delete(ctx context.Context, ID uint) error {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "BlogRepository.Delete"
	start := time.Now().UTC()

	_, err := r.db.Exec(ctx, blog_sql_Delete, ID)
	if err != nil {
		if strings.Contains(err.Error(), errMsg_duplicateKey) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return apperror.ErrBadRequest
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, blog_sql_Delete, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}
