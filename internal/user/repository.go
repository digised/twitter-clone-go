package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(
		ctx context.Context,
		user *User,
	) error

	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*User, error)

	GetByUsername(
		ctx context.Context,
		username string,
	) (*User, error)

	GetByEmail(
		ctx context.Context,
		email string,
	) (*User, error)

	UsernameExists(
		ctx context.Context,
		username string,
	) (bool, error)

	EmailExists(
		ctx context.Context,
		email string,
	) (bool, error)

	Update(
		ctx context.Context,
		user *User,
	) error

	UpdateProfile(
		ctx context.Context,
		id uuid.UUID,
		displayName *string,
		bio *string,
		avatarURL *string,
	) error

	Delete(
		ctx context.Context,
		id uuid.UUID,
	) error

	IncrementFollowersCount(
		ctx context.Context,
		id uuid.UUID,
	) error

	DecrementFollowersCount(
		ctx context.Context,
		id uuid.UUID,
	) error

	IncrementFollowingCount(
		ctx context.Context,
		id uuid.UUID,
	) error

	DecrementFollowingCount(
		ctx context.Context,
		id uuid.UUID,
	) error

	IncrementTweetsCount(
		ctx context.Context,
		id uuid.UUID,
	) error

	DecrementTweetsCount(
		ctx context.Context,
		id uuid.UUID,
	) error
}
