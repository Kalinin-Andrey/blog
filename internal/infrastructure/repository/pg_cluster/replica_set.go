package pg_cluster

import (
	"blog/internal/infrastructure/repository/pg"
)

type ReplicaSet struct {
	master *pg.Repository
	slave  *pg.Repository
}

func NewReplicaSet(dbMaster *pg.Repository, dbSlave *pg.Repository) *ReplicaSet {
	return &ReplicaSet{
		master: dbMaster,
		slave:  dbSlave,
	}
}

func (rs *ReplicaSet) WriteRepo() *pg.Repository {
	return rs.master
}

func (rs *ReplicaSet) ReadRepo() *pg.Repository {
	return rs.slave
}

func (rs *ReplicaSet) Close() {
	rs.master.Close()
	rs.slave.Close()
}
