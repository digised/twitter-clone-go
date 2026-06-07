package user

type CreateUserRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8,max=72"`
	DisplayName string `json:"display_name" binding:"required,min=1,max=100"`
}

type UpdateUserRequest struct {
	DisplayName *string `json:"display_name" binding:"omitempty,min=1,max=100"`
	Bio         *string `json:"bio" binding:"omitempty,max=160"`
	AvatarURL   *string `json:"avatar_url" binding:"omitempty,url"`
}

type UserResponse struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	DisplayName    string `json:"display_name"`
	Bio            string `json:"bio"`
	AvatarURL      string `json:"avatar_url"`
	FollowersCount int64  `json:"followers_count"`
	FollowingCount int64  `json:"following_count"`
	TweetsCount    int64  `json:"tweets_count"`
	IsVerified     bool   `json:"is_verified"`
}
