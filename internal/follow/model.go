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
