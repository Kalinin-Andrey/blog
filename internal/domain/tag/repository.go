package tag

import "context"

type ReplicaSet interface {
	WriteRepo() WriteRepository
	ReadRepo() ReadRepository
}

type WriteRepository interface {
	Create(ctx context.Context, entity *Tag) (ID uint, err error)
	Update(ctx context.Context, entity *Tag) error
	Delete(ctx context.Context, ID uint) error
}

type ReadRepository interface {
	Get(ctx context.Context, ID uint) (*Tag, error)
	MGet(ctx context.Context, IDs *[]uint) (*[]Tag, error)
	GetAll(ctx context.Context) (*[]Tag, error)
}
