package tweet

import (
	"time"

	"github.com/google/uuid"
)

type Tweet struct {
	ID           uuid.UUID  `db:"id" json:"id"`
	UserID       uuid.UUID  `db:"user_id" json:"user_id"`
	Content      string     `db:"content" json:"content"`
	MediaURLs    []string   `db:"media_urls" json:"media_urls,omitempty"`
	LikesCount   int        `db:"likes_count" json:"likes_count"`
	RetweetCount int        `db:"retweet_count" json:"retweet_count"`
	ReplyCount   int        `db:"reply_count" json:"reply_count"`
	ParentID     *uuid.UUID `db:"parent_id" json:"parent_id,omitempty"`
	RetweetedID  *uuid.UUID `db:"retweeted_id" json:"retweeted_id,omitempty"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
}

type CreateTweetDTO struct {
	Content   string     `json:"content" binding:"required,min=1,max=280"`
	MediaURLs []string   `json:"media_urls"`
	ParentID  *uuid.UUID `json:"parent_id"`
}

type UpdateTweetDTO struct {
	Content   *string   `json:"content" binding:"omitempty,min=1,max=280"`
	MediaURLs *[]string `json:"media_urls"`
}

type RetweetDTO struct {
	TweetID uuid.UUID `json:"tweet_id" binding:"required"`
}

type TweetResponse struct {
	ID             uuid.UUID      `json:"id"`
	Content        string         `json:"content"`
	MediaURLs      []string       `json:"media_urls,omitempty"`
	LikesCount     int            `json:"likes_count"`
	RetweetCount   int            `json:"retweet_count"`
	ReplyCount     int            `json:"reply_count"`
	ParentID       *uuid.UUID     `json:"parent_id,omitempty"`
	RetweetedID    *uuid.UUID     `json:"retweeted_id,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	User           UserResponse   `json:"user"`
	IsLiked        bool           `json:"is_liked"`
	IsRetweeted    bool           `json:"is_retweeted"`
	RetweetedTweet *TweetResponse `json:"retweeted_tweet,omitempty"`
}

type UserResponse struct {
	ID              uuid.UUID `json:"id"`
	Username        string    `json:"username"`
	DisplayName     string    `json:"display_name"`
	Bio             string    `json:"bio,omitempty"`
	ProfileImageURL string    `json:"profile_image_url,omitempty"`
	IsVerified      bool      `json:"is_verified"`
}
