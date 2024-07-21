package tsdb_cluster

import (
	"blog/internal/infrastructure/repository/tsdb"
)

type ReplicaSet struct {
	master *tsdb.Repository
	slave  *tsdb.Repository
}

func NewReplicaSet(dbMaster *tsdb.Repository, dbSlave *tsdb.Repository) *ReplicaSet {
	return &ReplicaSet{
		master: dbMaster,
		slave:  dbSlave,
	}
}

func (rs *ReplicaSet) WriteRepo() *tsdb.Repository {
	return rs.master
}

func (rs *ReplicaSet) ReadRepo() *tsdb.Repository {
	return rs.slave
}

func (rs *ReplicaSet) Close() {
	rs.master.Close()
	rs.slave.Close()
}
