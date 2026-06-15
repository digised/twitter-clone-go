package follow

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	IsFollowing(ctx context.Context, followerID, followedID uuid.UUID) (bool, error)
}

type Repository interface {
	Exists(ctx context.Context, followerID, followedID uuid.UUID) (bool, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	if repo == nil {
		panic("follow repository is nil")
	}
	return &service{repo: repo}
}

func (s *service) IsFollowing(ctx context.Context, followerID, followedID uuid.UUID) (bool, error) {
	return s.repo.Exists(ctx, followerID, followedID)
}
