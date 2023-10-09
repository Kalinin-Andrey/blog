package redis

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

func (rs *ReplicaSet) Master() *Repository {
	return rs.master
}

func (rs *ReplicaSet) Slave() *Repository {
	return rs.slave
}

func (rs *ReplicaSet) Close() error {
	err1 := rs.master.Close()
	err2 := rs.slave.Close()
	if err1 != nil {
		return err1
	}
	return err2
}
