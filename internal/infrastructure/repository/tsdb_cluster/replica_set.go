package tsdb_cluster

import (
	"github.com/Kalinin-Andrey/blog/internal/infrastructure/repository/tsdb"
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

func (rs *ReplicaSet) Master() *tsdb.Repository {
	return rs.master
}

func (rs *ReplicaSet) Slave() *tsdb.Repository {
	return rs.slave
}

func (rs *ReplicaSet) Close() {
	rs.master.Close()
	rs.slave.Close()
}
