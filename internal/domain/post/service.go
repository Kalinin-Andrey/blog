package post

import (
	"context"
	"github.com/minipkg/selection_condition"
)

type Service struct {
	replicaSet ReplicaSet
}

func NewService(replicaSet ReplicaSet) *Service {
	return &Service{
		replicaSet: replicaSet,
	}
}

func (s *Service) Create(ctx context.Context, entity *Post) (ID uint, err error) {
	return s.replicaSet.WriteRepo().Create(ctx, entity)
}

func (s *Service) Update(ctx context.Context, entity *Post) error {
	return s.replicaSet.WriteRepo().Update(ctx, entity)
}

func (s *Service) Delete(ctx context.Context, ID uint) error {
	return s.replicaSet.WriteRepo().Delete(ctx, ID)
}

func (s *Service) Get(ctx context.Context, ID uint) (*Post, error) {
	return s.replicaSet.ReadRepo().Get(ctx, ID)
}

func (s *Service) GetBySysname(ctx context.Context, sysname string) (*Post, error) {
	return s.replicaSet.ReadRepo().GetBySysname(ctx, sysname)
}

func (s *Service) Filter(ctx context.Context, blogSysname string, condition *selection_condition.SelectionCondition) (*[]PostPreview, error) {
	// todo
	return s.replicaSet.ReadRepo().Filter(ctx, condition)
}

func (s *Service) MGetByKeyword(ctx context.Context, blogSysname string, keywordSysname string, condition *selection_condition.SelectionCondition) (*[]PostPreview, error) {
	// todo
	return s.replicaSet.ReadRepo().Filter(ctx, condition)
}

func (s *Service) MGetByTag(ctx context.Context, blogSysname string, tagSysname string, condition *selection_condition.SelectionCondition) (*[]PostPreview, error) {
	// todo
	return s.replicaSet.ReadRepo().Filter(ctx, condition)
}

func (s *Service) TextSearch(ctx context.Context, searchString string, createdAtSortOrder *string) (*[]PostPreview, error) {
	// todo
	return s.replicaSet.ReadRepo().TextSearch(ctx, searchString, createdAtSortOrder)
}
