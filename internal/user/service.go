package user

import (
	"context"
	"log/slog"
	"time"
	"twitter-clone-go/pkg/metrics"

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
	log       *slog.Logger
}

func NewService(repo Repository, cache Cache, publisher EventPublisher, hasher PasswordHasher, log *slog.Logger) Service {
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
		log:       log,
	}
}

func (s *service) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
	usernameExists, err := s.repo.UsernameExists(ctx, req.Username)
	if err != nil {
		s.log.Error("failed to check username existence", "error", err, "username", req.Username)
		return nil, err
	}
	if usernameExists {
		return nil, ErrUsernameTaken
	}

	emailExists, err := s.repo.EmailExists(ctx, req.Email)
	if err != nil {
		s.log.Error("failed to check email existence", "error", err)
		return nil, err
	}
	if emailExists {
		return nil, ErrEmailTaken
	}

	passwordHash, err := s.hasher.Hash(req.Password)
	if err != nil {
		s.log.Error("failed to hash password", "error", err)
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
		s.log.Error("failed to create user", "error", err, "username", req.Username)
		return nil, err
	}

	s.log.Info("user created", "user_id", user.ID, "username", user.Username)
	metrics.UsersCreatedTotal.Inc()

	if s.publisher != nil {
		event := UserCreatedEvent{
			UserID:      user.ID,
			Username:    user.Username,
			Email:       user.Email,
			DisplayName: user.DisplayName,
			OccurredAt:  time.Now(),
		}
		if err := s.publisher.PublishUserCreated(ctx, event); err != nil {
			s.log.Warn("failed to publish user created event", "error", err, "user_id", user.ID)
		}
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
		s.log.Error("failed to get user by id", "error", err, "user_id", id)
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
		s.log.Error("failed to get user by username", "error", err, "username", username)
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
		s.log.Error("failed to update user profile", "error", err, "user_id", userID)
		return nil, err
	}

	s.log.Info("user profile updated", "user_id", userID)

	if s.cache != nil {
		_ = s.cache.DeleteUserByID(ctx, userID)
		_ = s.cache.DeleteUserByUsername(ctx, updatedUser.Username)
	}

	return updatedUser, nil
}

func (s *service) HandleCounterUpdate(ctx context.Context, userID uuid.UUID, fDelta, flDelta, tDelta int64) error {
	err := s.repo.UpdateCounters(ctx, userID, fDelta, flDelta, tDelta)
	if err != nil {
		s.log.Error("failed to update user counters", "error", err, "user_id", userID)
		return err
	}
	if s.cache != nil {
		_ = s.cache.DeleteUserByID(ctx, userID)
	}
	return nil
}
