package follow

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	FollowerID uuid.UUID `db:"follower_id"`
	FollowedID uuid.UUID `db:"followed_id"`
	CreatedAt  time.Time `db:"created_at"`
}

type FollowedUser struct {
	ID          uuid.UUID `db:"id"`
	Username    string    `db:"username"`
	DisplayName string    `db:"display_name"`
	AvatarURL   string    `db:"avatar_url"`
	IsVerified  bool      `db:"is_verified"`
}
