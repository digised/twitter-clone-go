package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserCreatedEvent struct {
	UserID      uuid.UUID `json:"user_id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	OccurredAt  time.Time `json:"occurred_at"`
}

type EventPublisher interface {
	PublishUserCreated(ctx context.Context, event UserCreatedEvent) error
}
