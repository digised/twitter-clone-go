package application

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"twitter-clone-go/internal/follow"
	"twitter-clone-go/internal/user"
	"twitter-clone-go/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListFollowHandler struct {
	userService   user.Service
	followService follow.Service
	log           *slog.Logger
}

func NewListFollowHandler(u user.Service, f follow.Service, log *slog.Logger) *ListFollowHandler {
	if u == nil {
		panic("user service is nil")
	}
	if f == nil {
		panic("follow service is nil")
	}
	return &ListFollowHandler{userService: u, followService: f, log: log}
}

func (h *ListFollowHandler) GetFollowers(c *gin.Context) {
	h.list(c, h.followService.GetFollowers)
}

func (h *ListFollowHandler) GetFollowing(c *gin.Context) {
	h.list(c, h.followService.GetFollowing)
}

func (h *ListFollowHandler) list(
	c *gin.Context,
	fn func(ctx context.Context, userID uuid.UUID, p utils.Pagination) (utils.PaginatedResponse[follow.FollowUserResponse], error),
) {
	username := c.Param("username")

	target, err := h.userService.GetByUsername(c.Request.Context(), username)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			utils.RespondError(c, http.StatusNotFound, "user not found")
			return
		}
		h.log.Error("failed to resolve user for follow list", "error", err, "username", username)
		utils.RespondError(c, http.StatusInternalServerError, "internal server error")
		return
	}

	p := utils.ParsePagination(c)

	res, err := fn(c.Request.Context(), target.ID, p)
	if err != nil {
		h.log.Error("failed to list follows", "error", err, "user_id", target.ID)
		utils.RespondError(c, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.RespondOK(c, res)
}
