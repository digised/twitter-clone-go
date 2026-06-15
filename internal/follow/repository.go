package follow

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Exists(ctx context.Context, followerID, followedID uuid.UUID) (bool, error)
	Create(ctx context.Context, followerID, followedID uuid.UUID) error
	Delete(ctx context.Context, followerID, followedID uuid.UUID) error
	GetFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]FollowedUser, int64, error)
	GetFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) ([]FollowedUser, int64, error)
	CountFollowsSince(ctx context.Context, followerID uuid.UUID, since time.Time) (int64, error)
}
