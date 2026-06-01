package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Cache interface {
	SetUserByID(
		ctx context.Context,
		user *User,
		ttl time.Duration,
	) error
	SetUserByUsername(
		ctx context.Context,
		user *User,
		ttl time.Duration,
	) error

	GetUserByID(
		ctx context.Context,
		id uuid.UUID,
	) (*User, error)
	GetUserByUsername(
		ctx context.Context,
		username string,
	) (*User, error)

	DeleteUserByID(
		ctx context.Context,
		id uuid.UUID,
	) error
	DeleteUserByUsername(
		ctx context.Context,
		username string,
	) error
}
