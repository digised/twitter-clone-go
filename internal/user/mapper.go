package user

func ToUserResponse(user User) UserResponse {
	return UserResponse{
		ID:             user.ID.String(),
		Username:       user.Username,
		DisplayName:    user.DisplayName,
		Bio:            user.Bio,
		AvatarURL:      user.AvatarURL,
		FollowersCount: user.FollowersCount,
		FollowingCount: user.FollowingCount,
		TweetsCount:    user.TweetsCount,
		IsVerified:     user.IsVerified,
	}
}
