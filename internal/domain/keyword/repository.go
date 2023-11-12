package keyword

import "context"

type ReplicaSet interface {
	WriteRepo() WriteRepository
	ReadRepo() ReadRepository
}

type WriteRepository interface {
	Create(ctx context.Context, entity *Keyword) (ID uint, err error)
	Update(ctx context.Context, entity *Keyword) error
	Delete(ctx context.Context, ID uint) error
}

type ReadRepository interface {
	Get(ctx context.Context, ID uint) (*Keyword, error)
	MGet(ctx context.Context, ID *[]uint) (*[]Keyword, error)
	GetAll(ctx context.Context) (*[]Keyword, error)
}
