package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const ContextUserIDKey = "user_id"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	if service == nil {
		panic("user service is nil")
	}
	return &Handler{service: service}
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
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.CreateUser(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrUsernameTaken) || errors.Is(err, ErrEmailTaken) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, ToUserResponse(*user))
}

func (h *Handler) GetProfile(c *gin.Context) {
	username := c.Param("username")

	user, err := h.service.GetByUsername(c.Request.Context(), username)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, ToUserResponse(*user))
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	var req UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	userIDVal, exists := c.Get(ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user identity"})
		return
	}

	user, err := h.service.UpdateUser(c.Request.Context(), userID, req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, ToUserResponse(*user))
}
