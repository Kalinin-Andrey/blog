package redis

import (
	"errors"
)

type ReplicaSet struct {
	master *Repository
	slave  *Repository
}

func NewReplicaSet(dbMaster *Repository, dbSlave *Repository) *ReplicaSet {
	return &ReplicaSet{
		master: dbMaster,
		slave:  dbSlave,
	}
}

func (rs *ReplicaSet) WriteRepo() *Repository {
	return rs.master
}

func (rs *ReplicaSet) ReadRepo() *Repository {
	return rs.slave
}

func (rs *ReplicaSet) Close() error {
	return errors.Join(rs.master.Close(), rs.slave.Close())
}
