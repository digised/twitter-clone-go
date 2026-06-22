package user

import (
	"errors"
	"log/slog"
	"net/http"

	"twitter-clone-go/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const ContextUserIDKey = "user_id"

type Handler struct {
	service Service
	log     *slog.Logger
}

func NewHandler(service Service, log *slog.Logger) *Handler {
	if service == nil {
		panic("user service is nil")
	}
	return &Handler{service: service, log: log}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	public := rg.Group("/users")
	{
		public.POST("", h.RegisterUser)
		public.GET("/:username", h.GetProfile)
	}

	protected := rg.Group("/users")
	protected.Use(authMiddleware)
	{
		protected.PATCH("/me", h.UpdateProfile)
	}
}

func (h *Handler) RegisterUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user, err := h.service.CreateUser(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrUsernameTaken) || errors.Is(err, ErrEmailTaken) {
			utils.RespondError(c, http.StatusConflict, err.Error())
			return
		}
		h.log.Error("unexpected error creating user", "error", err)
		utils.RespondError(c, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.RespondCreated(c, ToUserResponse(*user))
}

func (h *Handler) GetProfile(c *gin.Context) {
	username := c.Param("username")

	user, err := h.service.GetByUsername(c.Request.Context(), username)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			utils.RespondError(c, http.StatusNotFound, "user not found")
			return
		}
		h.log.Error("unexpected error fetching user profile", "error", err, "username", username)
		utils.RespondError(c, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.RespondOK(c, ToUserResponse(*user))
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	var req UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	userIDVal, exists := c.Get(ContextUserIDKey)
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		h.log.Error("invalid user_id type in context", "value", userIDVal)
		utils.RespondError(c, http.StatusInternalServerError, "invalid user identity")
		return
	}

	user, err := h.service.UpdateUser(c.Request.Context(), userID, req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			utils.RespondError(c, http.StatusNotFound, "user not found")
			return
		}
		h.log.Error("unexpected error updating user profile", "error", err, "user_id", userID)
		utils.RespondError(c, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.RespondOK(c, ToUserResponse(*user))
}
