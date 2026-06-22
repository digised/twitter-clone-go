package follow

import (
	"context"
	"log/slog"
	"time"
	"twitter-clone-go/internal/constants"
	"twitter-clone-go/internal/utils"
	"twitter-clone-go/pkg/metrics"

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
	log       *slog.Logger
}

func NewService(repo Repository, publisher EventPublisher, log *slog.Logger) Service {
	if repo == nil {
		panic("follow repository is nil")
	}
	return &service{repo: repo, publisher: publisher, log: log}
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
		s.log.Error("failed to count follows", "error", err, "follower_id", followerID)
		return err
	}
	if count >= constants.MaxFollowsPerHour {
		s.log.Warn("follow rate limit exceeded", "follower_id", followerID, "count", count)
		metrics.FollowRateLimitHitsTotal.Inc()
		return ErrFollowRateLimitExceeded
	}

	exists, err := s.repo.Exists(ctx, followerID, followedID)
	if err != nil {
		s.log.Error("failed to check follow existence", "error", err, "follower_id", followerID, "followed_id", followedID)
		return err
	}
	if exists {
		return ErrAlreadyFollowing
	}

	if err := s.repo.Create(ctx, followerID, followedID); err != nil {
		s.log.Error("failed to create follow", "error", err, "follower_id", followerID, "followed_id", followedID)
		return err
	}

	s.log.Info("user followed", "follower_id", followerID, "followed_id", followedID)
	metrics.FollowsTotal.Inc()

	if s.publisher != nil {
		event := UserFollowedEvent{
			FollowerID: followerID,
			FollowedID: followedID,
			OccurredAt: time.Now(),
		}
		if err := s.publisher.PublishUserFollowed(ctx, event); err != nil {
			s.log.Warn("failed to publish user followed event", "error", err, "follower_id", followerID)
		}
	}

	return nil
}

func (s *service) Unfollow(ctx context.Context, followerID, followedID uuid.UUID) error {
	exists, err := s.repo.Exists(ctx, followerID, followedID)
	if err != nil {
		s.log.Error("failed to check follow existence", "error", err, "follower_id", followerID, "followed_id", followedID)
		return err
	}
	if !exists {
		return ErrNotFollowing
	}

	if err := s.repo.Delete(ctx, followerID, followedID); err != nil {
		s.log.Error("failed to delete follow", "error", err, "follower_id", followerID, "followed_id", followedID)
		return err
	}

	s.log.Info("user unfollowed", "follower_id", followerID, "followed_id", followedID)
	metrics.UnfollowsTotal.Inc()

	return nil
}

func (s *service) GetFollowers(ctx context.Context, userID uuid.UUID, p utils.Pagination) (utils.PaginatedResponse[FollowUserResponse], error) {
	users, total, err := s.repo.GetFollowers(ctx, userID, p.Limit, p.Offset)
	if err != nil {
		s.log.Error("failed to get followers", "error", err, "user_id", userID)
		return utils.PaginatedResponse[FollowUserResponse]{}, err
	}
	return utils.NewPaginatedResponse(toFollowUserResponses(users), total, p), nil
}

func (s *service) GetFollowing(ctx context.Context, userID uuid.UUID, p utils.Pagination) (utils.PaginatedResponse[FollowUserResponse], error) {
	users, total, err := s.repo.GetFollowing(ctx, userID, p.Limit, p.Offset)
	if err != nil {
		s.log.Error("failed to get following", "error", err, "user_id", userID)
		return utils.PaginatedResponse[FollowUserResponse]{}, err
	}
	return utils.NewPaginatedResponse(toFollowUserResponses(users), total, p), nil
}
