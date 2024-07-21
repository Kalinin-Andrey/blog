package blog

import (
	"context"
)

type ReplicaSet interface {
	WriteRepo() WriteRepository
	ReadRepo() ReadRepository
}

type WriteRepository interface {
	Create(ctx context.Context, entity *Blog) (ID uint, err error)
	Update(ctx context.Context, entity *Blog) error
	Delete(ctx context.Context, ID uint) error
}

type ReadRepository interface {
	Get(ctx context.Context, sysname string) (*Blog, error)
	MGet(ctx context.Context, sysnames *[]string) (*[]Blog, error)
	GetAll(ctx context.Context) (*[]Blog, error)
}
