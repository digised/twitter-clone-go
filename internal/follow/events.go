package follow

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserFollowedEvent struct {
	FollowerID uuid.UUID `json:"follower_id"`
	FollowedID uuid.UUID `json:"followed_id"`
	OccurredAt time.Time `json:"occurred_at"`
}

type EventPublisher interface {
	PublishUserFollowed(ctx context.Context, event UserFollowedEvent) error
}
