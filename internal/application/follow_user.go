package application

import (
	"errors"
	"net/http"

	"twitter-clone-go/internal/follow"
	"twitter-clone-go/internal/user"
	"twitter-clone-go/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FollowUserHandler struct {
	userService   user.Service
	followService follow.Service
}

func NewFollowUserHandler(u user.Service, f follow.Service) *FollowUserHandler {
	if u == nil {
		panic("user service is nil")
	}
	if f == nil {
		panic("follow service is nil")
	}
	return &FollowUserHandler{userService: u, followService: f}
}

func (h *FollowUserHandler) Follow(c *gin.Context) {
	followerID, target, ok := h.resolve(c)
	if !ok {
		return
	}

	if err := h.followService.Follow(c.Request.Context(), followerID, target.ID); err != nil {
		writeFollowError(c, err)
		return
	}

	_ = h.userService.HandleCounterUpdate(c.Request.Context(), followerID, 0, 1, 0)
	_ = h.userService.HandleCounterUpdate(c.Request.Context(), target.ID, 1, 0, 0)

	utils.RespondStatusOK(c)
}

func (h *FollowUserHandler) Unfollow(c *gin.Context) {
	followerID, target, ok := h.resolve(c)
	if !ok {
		return
	}

	if err := h.followService.Unfollow(c.Request.Context(), followerID, target.ID); err != nil {
		writeFollowError(c, err)
		return
	}

	_ = h.userService.HandleCounterUpdate(c.Request.Context(), followerID, 0, -1, 0)
	_ = h.userService.HandleCounterUpdate(c.Request.Context(), target.ID, -1, 0, 0)

	utils.RespondStatusOK(c)
}

func (h *FollowUserHandler) resolve(c *gin.Context) (followerID uuid.UUID, target *user.User, ok bool) {
	userIDVal, exists := c.Get(user.ContextUserIDKey)
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return uuid.UUID{}, nil, false
	}

	followerID, valid := userIDVal.(uuid.UUID)
	if !valid {
		utils.RespondError(c, http.StatusInternalServerError, "invalid user identity")
		return uuid.UUID{}, nil, false
	}

	username := c.Param("username")

	target, err := h.userService.GetByUsername(c.Request.Context(), username)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			utils.RespondError(c, http.StatusNotFound, "user not found")
			return uuid.UUID{}, nil, false
		}
		utils.RespondError(c, http.StatusInternalServerError, "internal server error")
		return uuid.UUID{}, nil, false
	}

	return followerID, target, true
}

func writeFollowError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, follow.ErrCannotFollowSelf):
		utils.RespondError(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, follow.ErrAlreadyFollowing), errors.Is(err, follow.ErrNotFollowing):
		utils.RespondError(c, http.StatusConflict, err.Error())
	case errors.Is(err, follow.ErrFollowRateLimitExceeded):
		utils.RespondError(c, http.StatusTooManyRequests, err.Error())
	default:
		utils.RespondError(c, http.StatusInternalServerError, "internal server error")
	}
}
