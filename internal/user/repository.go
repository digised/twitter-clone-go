package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, displayName *string, bio *string, avatarURL *string) (*User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UsernameExists(ctx context.Context, username string) (bool, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	UpdateCounters(ctx context.Context, id uuid.UUID, followersDelta, followingDelta, tweetsDelta int64) error
}
