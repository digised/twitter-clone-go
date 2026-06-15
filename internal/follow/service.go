package follow

import (
	"context"
	"time"
	"twitter-clone-go/internal/constants"

	"twitter-clone-go/internal/utils"

	"github.com/google/uuid"
)

type Service interface {
	IsFollowing(ctx context.Context, followerID, followedID uuid.UUID) (bool, error)
	Follow(ctx context.Context, followerID, followedID uuid.UUID) error
	Unfollow(ctx context.Context, followerID, followedID uuid.UUID) error
	GetFollowers(ctx context.Context, userID uuid.UUID, p utils.Pagination) (utils.PaginatedResponse[FollowUserResponse], error)
	GetFollowing(ctx context.Context, userID uuid.UUID, p utils.Pagination) (utils.PaginatedResponse[FollowUserResponse], error)
}

type service struct {
	repo      Repository
	publisher EventPublisher
}

func NewService(repo Repository, publisher EventPublisher) Service {
	if repo == nil {
		panic("follow repository is nil")
	}
	return &service{repo: repo, publisher: publisher}
}

func (s *service) IsFollowing(ctx context.Context, followerID, followedID uuid.UUID) (bool, error) {
	return s.repo.Exists(ctx, followerID, followedID)
}

func (s *service) Follow(ctx context.Context, followerID, followedID uuid.UUID) error {
	if followerID == followedID {
		return ErrCannotFollowSelf
	}

	count, err := s.repo.CountFollowsSince(ctx, followerID, time.Now().Add(-1*time.Hour))
	if err != nil {
		return err
	}
	if count >= constants.MaxFollowsPerHour {
		return ErrFollowRateLimitExceeded
	}

	exists, err := s.repo.Exists(ctx, followerID, followedID)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyFollowing
	}

	if err := s.repo.Create(ctx, followerID, followedID); err != nil {
		return err
	}

	if s.publisher != nil {
		event := UserFollowedEvent{
			FollowerID: followerID,
			FollowedID: followedID,
			OccurredAt: time.Now(),
		}
		_ = s.publisher.PublishUserFollowed(ctx, event)
	}

	return nil
}

func (s *service) Unfollow(ctx context.Context, followerID, followedID uuid.UUID) error {
	exists, err := s.repo.Exists(ctx, followerID, followedID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotFollowing
	}

	return s.repo.Delete(ctx, followerID, followedID)
}

func (s *service) GetFollowers(ctx context.Context, userID uuid.UUID, p utils.Pagination) (utils.PaginatedResponse[FollowUserResponse], error) {
	users, total, err := s.repo.GetFollowers(ctx, userID, p.Limit, p.Offset)
	if err != nil {
		return utils.PaginatedResponse[FollowUserResponse]{}, err
	}
	return utils.NewPaginatedResponse(toFollowUserResponses(users), total, p), nil
}

func (s *service) GetFollowing(ctx context.Context, userID uuid.UUID, p utils.Pagination) (utils.PaginatedResponse[FollowUserResponse], error) {
	users, total, err := s.repo.GetFollowing(ctx, userID, p.Limit, p.Offset)
	if err != nil {
		return utils.PaginatedResponse[FollowUserResponse]{}, err
	}
	return utils.NewPaginatedResponse(toFollowUserResponses(users), total, p), nil
}
