package application

import (
	"log/slog"
	"net/http"

	"twitter-clone-go/internal/follow"
	"twitter-clone-go/internal/user"
	"twitter-clone-go/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserProfileQueryHandler struct {
	userService   user.Service
	followService follow.Service
	log           *slog.Logger
}

type UserProfileResponse struct {
	user.UserResponse
	IsFollowing  bool `json:"is_following"`
	IsOwnProfile bool `json:"is_own_profile"`
}

func NewUserProfileQueryHandler(u user.Service, f follow.Service, log *slog.Logger) *UserProfileQueryHandler {
	if u == nil {
		panic("user service is nil")
	}
	if f == nil {
		panic("follow service is nil")
	}
	return &UserProfileQueryHandler{userService: u, followService: f, log: log}
}

func (h *UserProfileQueryHandler) Execute(c *gin.Context) {
	username := c.Param("username")

	u, err := h.userService.GetByUsername(c.Request.Context(), username)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, "user not found")
		return
	}

	var isFollowing, isOwnProfile bool
	userIDVal, exists := c.Get(user.ContextUserIDKey)

	if exists {
		if viewerID, ok := userIDVal.(uuid.UUID); ok {
			isOwnProfile = viewerID == u.ID
			if !isOwnProfile {
				var err error
				isFollowing, err = h.followService.IsFollowing(c.Request.Context(), viewerID, u.ID)
				if err != nil {
					h.log.Warn("failed to check isFollowing", "error", err, "viewer_id", viewerID, "target_id", u.ID)
				}
			}
		}
	}

	response := UserProfileResponse{
		UserResponse: user.ToUserResponse(*u),
		IsFollowing:  isFollowing,
		IsOwnProfile: isOwnProfile,
	}

	utils.RespondOK(c, response)
}
