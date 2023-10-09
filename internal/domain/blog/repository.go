package blog

import (
	"context"
)

type TsDBCluster interface {
	GetShardMaster(num byte) TsDBRepository
	GetShardSlave(num byte) TsDBRepository
	GetShardMasterByID(ID uint) TsDBRepository
	GetShardSlaveByID(ID uint) TsDBRepository
}

type TsDBReplicaSet interface {
	Master() TsDBRepository
	Slave() TsDBRepository
}

type TsDBRepository interface {
	Create(ctx context.Context, entity *Blog) error
	MGet(ctx context.Context, filter *Filter4TsDB) (*[]Blog, error)
}

type FastCluster interface {
	GetShardMasterByID(ID uint) FastRepository
	GetShardSlaveByID(ID uint) FastRepository
	MCreate(ctx context.Context, entities *[]Blog) error
	MGet(ctx context.Context, sellerIDs *[]uint) (*[]Blog, error)
}

type FastReplicaSet interface {
	Master() FastRepository
	Slave() FastRepository
}

type FastRepository interface {
	Get(ctx context.Context, ID uint) (*Blog, error)
	MGet(ctx context.Context, IDs *[]uint) (*[]Blog, error)
	Set(ctx context.Context, blog *Blog) error
	MSet(ctx context.Context, blogList *[]Blog) error
}
