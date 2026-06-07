package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type Service interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, req UpdateUserRequest) (*User, error)
	HandleCounterUpdate(ctx context.Context, userID uuid.UUID, fDelta, flDelta, tDelta int64) error
}

type service struct {
	repo      Repository
	cache     Cache
	publisher EventPublisher
	hasher    PasswordHasher
}

func NewService(repo Repository, cache Cache, publisher EventPublisher, hasher PasswordHasher) Service {
	if repo == nil {
		panic("user repository is nil")
	}
	if hasher == nil {
		panic("password hasher is nil")
	}
	return &service{
		repo:      repo,
		cache:     cache,
		publisher: publisher,
		hasher:    hasher,
	}
}

func (s *service) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
	usernameExists, err := s.repo.UsernameExists(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if usernameExists {
		return nil, ErrUsernameTaken
	}

	emailExists, err := s.repo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, ErrEmailTaken
	}

	passwordHash, err := s.hasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:           uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		DisplayName:  req.DisplayName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	if s.publisher != nil {
		event := UserCreatedEvent{
			UserID:      user.ID,
			Username:    user.Username,
			Email:       user.Email,
			DisplayName: user.DisplayName,
			OccurredAt:  time.Now(),
		}
		_ = s.publisher.PublishUserCreated(ctx, event)
	}

	return user, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	if s.cache != nil {
		if cached, err := s.cache.GetUserByID(ctx, id); err == nil {
			return cached, nil
		}
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}

	if s.cache != nil {
		_ = s.cache.SetUserByID(ctx, user, 10*time.Minute)
	}

	return user, nil
}

func (s *service) GetByUsername(ctx context.Context, username string) (*User, error) {
	if s.cache != nil {
		if cached, err := s.cache.GetUserByUsername(ctx, username); err == nil {
			return cached, nil
		}
	}

	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, ErrNotFound
	}

	if s.cache != nil {
		_ = s.cache.SetUserByUsername(ctx, user, 10*time.Minute)
	}

	return user, nil
}

func (s *service) UpdateUser(ctx context.Context, userID uuid.UUID, req UpdateUserRequest) (*User, error) {
	updatedUser, err := s.repo.UpdateProfile(ctx, userID, req.DisplayName, req.Bio, req.AvatarURL)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		_ = s.cache.DeleteUserByID(ctx, userID)
		_ = s.cache.DeleteUserByUsername(ctx, updatedUser.Username)
	}

	return updatedUser, nil
}

func (s *service) HandleCounterUpdate(ctx context.Context, userID uuid.UUID, fDelta, flDelta, tDelta int64) error {
	err := s.repo.UpdateCounters(ctx, userID, fDelta, flDelta, tDelta)
	if err != nil {
		return err
	}
	if s.cache != nil {
		_ = s.cache.DeleteUserByID(ctx, userID)
	}
	return nil
}
