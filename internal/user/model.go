package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `db:"id"`
	Username       string    `db:"username"`
	Email          string    `db:"email"`
	PasswordHash   string    `db:"password_hash"`
	DisplayName    string    `db:"display_name"`
	Bio            string    `db:"bio"`
	AvatarURL      string    `db:"avatar_url"`
	FollowersCount int64     `db:"followers_count"`
	FollowingCount int64     `db:"following_count"`
	TweetsCount    int64     `db:"tweets_count"`
	IsVerified     bool      `db:"is_verified"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
