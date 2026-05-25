package user

func ToUserResponse(user User) UserResponse {
	return UserResponse{
		ID:             user.ID,
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

func ToUserProfileResponse(
	user User,
	isFollowing bool,
	isOwnProfile bool,
) UserProfileResponse {
	return UserProfileResponse{
		ID:             user.ID,
		Username:       user.Username,
		DisplayName:    user.DisplayName,
		Bio:            user.Bio,
		AvatarURL:      user.AvatarURL,
		FollowersCount: user.FollowersCount,
		FollowingCount: user.FollowingCount,
		TweetsCount:    user.TweetsCount,
		IsVerified:     user.IsVerified,
		IsFollowing:    isFollowing,
		IsOwnProfile:   isOwnProfile,
	}
}
