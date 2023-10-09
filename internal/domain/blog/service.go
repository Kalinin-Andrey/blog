package blog

type Service struct {
	fastReplicaSet FastReplicaSet
	tsDBReplicaSet TsDBReplicaSet
}

func NewService(fastReplicaSet FastReplicaSet, tsDBReplicaSet TsDBReplicaSet) *Service {
	return &Service{
		fastReplicaSet: fastReplicaSet,
		tsDBReplicaSet: tsDBReplicaSet,
	}
}

/**
func (s *Service) Get(ctx context.Context, sellerID uint) (*Blog, error) {
	return s.fastReplicaSet.GetShardSlaveBySellerID(sellerID).Get(ctx, sellerID)
}

func (s *Service) MGet(ctx context.Context, sellerIDs *[]uint) (*[]Blog, error) {
	return s.fastReplicaSet.MGet(ctx, sellerIDs)
}

func (s *Service) Filter(ctx context.Context, filter *Filter) (count uint, list *[]Blog, err error) {
	count, err = s.replicaSet.Slave().Count(ctx, filter)
	if err != nil {
		return 0, nil, err
	}
	if count == 0 {
		return 0, nil, nil
	}

	list, err = s.replicaSet.Slave().MGet(ctx, filter)
	if err != nil {
		return 0, nil, err
	}
	rs := Ratings(*list)
	sellerOldIDs := rs.GetSellerOldIDs()
	listR, err := s.MGet(ctx, sellerOldIDs) // делаем выборку из Redis, чтобы взять из неё SellerName
	if err != nil {
		return 0, nil, err
	}

	listByOldSellerID := make(map[uint]Blog, len(*listR))
	for _, item := range *listR {
		listByOldSellerID[item.SellerOldId] = item
	}

	for i := range *list {
		item, ok := listByOldSellerID[(*list)[i].SellerOldId]
		if !ok {
			continue
		}
		(*list)[i].SellerName = item.SellerName
	}
	return count, list, nil
}

func (s *Service) Create(ctx context.Context, entity *Blog) error {
	err := s.replicaSet.Master().Create(ctx, entity)
	if err != nil {
		return err
	}
	return s.fastReplicaSet.GetShardMasterBySellerID(entity.SellerOldId).Set(ctx, entity)
}

func (s *Service) MCreate(ctx context.Context, entities *[]Blog) (err error) {
	if err = s.replicaSet.Master().MCreate(ctx, entities); err != nil {
		return err
	}

	return s.fastReplicaSet.MCreate(ctx, entities)
}
**/
/**
* 3 funcs for storing resulting rating in TsDB
func (s *Service) FilterInTsDB(ctx context.Context, filter *Filter4TsDB) (*[]Blog, error) {
	list, err := s.tsDBReplicaSet.MGet(ctx, filter)
	if err != nil {
		return nil, err
	}
	rs := Ratings(*list)
	sellerOldIDs := rs.GetSellerOldIDs()
	listR, err := s.MGet(ctx, sellerOldIDs)
	if err != nil {
		return nil, err
	}

	listByOldSellerID := make(map[uint]Blog, len(*listR))
	for _, item := range *listR {
		listByOldSellerID[item.SellerOldId] = item
	}

	for i := range *list {
		item, ok := listByOldSellerID[(*list)[i].SellerOldId]
		if !ok {
			continue
		}
		(*list)[i].SellerName = item.SellerName
	}
	return list, nil
}

func (s *Service) CreateInTsDB(ctx context.Context, entity *Blog) error {
	err := s.tsDBReplicaSet.GetShardMasterBySellerID(entity.SellerOldId).Create(ctx, entity)
	if err != nil {
		return err
	}
	//return s.fastReplicaSet.GetShardMasterBySellerID(entity.SellerOldId).Set(ctx, entity)
	return nil
}

func (s *Service) MCreateInTsDB(ctx context.Context, entities *[]Blog) (err error) {
	if err = s.tsDBReplicaSet.MCreate(ctx, entities); err != nil {
		return err
	}
	// todo:+
	//return s.fastReplicaSet.MCreate(ctx, entities)
	return nil
}
*/
