package blog

import (
	"context"
)

type TsDBCluster interface {
	GetShardWriteRepo(num byte) WriteTsDBRepository
	GetShardReadRepo(num byte) ReadTsDBRepository
	GetShardWriteRepoByID(ID uint) WriteTsDBRepository
	GetShardReadRepoByID(ID uint) ReadTsDBRepository
}

type TsDBReplicaSet interface {
	WriteRepo() WriteTsDBRepository
	ReadRepo() ReadTsDBRepository
}

type WriteTsDBRepository interface {
	Create(ctx context.Context, entity *Blog) error
}

type ReadTsDBRepository interface {
	MGet(ctx context.Context, filter *Filter4TsDB) (*[]Blog, error)
}

type FastCluster interface {
	GetShardWriteRepoByID(ID uint) WriteFastRepository
	GetShardReadRepoByID(ID uint) ReadFastRepository
}

type FastReplicaSet interface {
	WriteRepo() WriteFastRepository
	ReadRepo() ReadFastRepository
}

type WriteFastRepository interface {
	Set(ctx context.Context, blog *Blog) error
	MSet(ctx context.Context, blogList *[]Blog) error
}

type ReadFastRepository interface {
	Get(ctx context.Context, ID uint) (*Blog, error)
	MGet(ctx context.Context, IDs *[]uint) (*[]Blog, error)
}
