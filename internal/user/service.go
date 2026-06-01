package user

import (
	"context"

	"github.com/google/uuid"
)

type service struct {
	repo  Repository
	cache Cache
}

func NewService(
	repo Repository,
	cache Cache,
) Service {
	if repo == nil {
		panic("repository is null")
	}

	return &service{
		repo:  repo,
		cache: cache,
	}
}

type Service interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*User, error)
	Login(ctx context.Context, req LoginRequest) (string, *User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, req UpdateUserRequest) (*User, error)
	IsFollowing(ctx context.Context, followerID uuid.UUID, followingID uuid.UUID) (bool, error)
}

func (s *service) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
	panic("null")
}

func (s *service) Login(ctx context.Context, req LoginRequest) (string, *User, error) {
	panic("null")
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	panic("null")
}

func (s *service) GetByUsername(ctx context.Context, username string) (*User, error) {
	panic("null")
}

func (s *service) UpdateUser(ctx context.Context, userID uuid.UUID, req UpdateUserRequest) (*User, error) {
	panic("null")
}

func (s *service) IsFollowing(ctx context.Context, followerID uuid.UUID, followingID uuid.UUID) (bool, error) {
	panic("null")
}
