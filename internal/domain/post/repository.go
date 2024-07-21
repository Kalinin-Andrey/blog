package post

import (
	"context"
	"github.com/minipkg/selection_condition"
)

type ReplicaSet interface {
	WriteRepo() WriteRepository
	ReadRepo() ReadRepository
}

type WriteRepository interface {
	Create(ctx context.Context, entity *Post) (ID uint, err error)
	Update(ctx context.Context, entity *Post) error
	Delete(ctx context.Context, ID uint) error
}

type ReadRepository interface {
	Get(ctx context.Context, ID uint) (*Post, error)
	GetBySysname(ctx context.Context, sysname string) (*Post, error)
	MGet(ctx context.Context, IDs *[]uint) (*[]PostPreview, error)
	Filter(ctx context.Context, condition *selection_condition.SelectionCondition) (*[]PostPreview, error)
	TextSearch(ctx context.Context, searchString string, createdAtSortOrder *string) (*[]PostPreview, error)
}
