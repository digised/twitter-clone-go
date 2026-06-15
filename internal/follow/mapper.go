package follow

func ToFollowUserResponse(u FollowedUser) FollowUserResponse { // same as FollowedUser
	return FollowUserResponse{
		ID:          u.ID.String(),
		Username:    u.Username,
		DisplayName: u.DisplayName,
		AvatarURL:   u.AvatarURL,
		IsVerified:  u.IsVerified,
	}
}

func toFollowUserResponses(users []FollowedUser) []FollowUserResponse {
	res := make([]FollowUserResponse, 0, len(users))
	for _, u := range users {
		res = append(res, ToFollowUserResponse(u))
	}
	return res
}
