package application

import (
	_ "context"
	"net/http"

	"twitter-clone-go/internal/follow"
	"twitter-clone-go/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserProfileQueryHandler struct {
	userService   user.Service
	followService follow.Service
}

type UserProfileResponse struct {
	user.UserResponse
	IsFollowing  bool `json:"is_following"`
	IsOwnProfile bool `json:"is_own_profile"`
}

func NewUserProfileQueryHandler(u user.Service, f follow.Service) *UserProfileQueryHandler {
	if u == nil {
		panic("user service is nil")
	}
	if f == nil {
		panic("follow service is nil")
	}
	return &UserProfileQueryHandler{userService: u, followService: f}
}

func (h *UserProfileQueryHandler) Execute(c *gin.Context) {
	username := c.Param("username")

	u, err := h.userService.GetByUsername(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var isFollowing, isOwnProfile bool
	userIDVal, exists := c.Get(user.ContextUserIDKey)

	if exists {
		if viewerID, ok := userIDVal.(uuid.UUID); ok {
			isOwnProfile = viewerID == u.ID
			if !isOwnProfile {
				isFollowing, _ = h.followService.IsFollowing(c.Request.Context(), viewerID, u.ID)
			}
		}
	}

	response := UserProfileResponse{
		UserResponse: user.ToUserResponse(*u),
		IsFollowing:  isFollowing,
		IsOwnProfile: isOwnProfile,
	}

	c.JSON(http.StatusOK, response)
}
