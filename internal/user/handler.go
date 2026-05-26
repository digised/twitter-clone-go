package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Service interface {
	CreateUser(
		ctx context.Context,
		req CreateUserRequest,
	) (*User, error)

	Login(
		ctx context.Context,
		req LoginRequest,
	) (string, *User, error)

	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*User, error)

	GetByUsername(
		ctx context.Context,
		username string,
	) (*User, error)

	UpdateUser(
		ctx context.Context,
		userID uuid.UUID,
		req UpdateUserRequest,
	) (*User, error)

	IsFollowing(
		ctx context.Context,
		followerID uuid.UUID,
		followingID uuid.UUID,
	) (bool, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		users.GET("/:username", h.GetProfile)
		users.PATCH("/me", h.UpdateProfile)
	}
	auth := rg.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.service.CreateUser(
		c.Request.Context(),
		req,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(
		http.StatusCreated,
		ToUserResponse(*user),
	)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, user, err := h.service.Login(
		c.Request.Context(),
		req,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  ToUserResponse(*user),
	})
}

func (h *Handler) GetProfile(c *gin.Context) {
	username := c.Param("username")

	user, err := h.service.GetByUsername(
		c.Request.Context(),
		username,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	var (
		isFollowing  bool
		isOwnProfile bool
	)

	userIDValue, exists := c.Get("user_id")
	if exists {
		viewerID, ok := userIDValue.(uuid.UUID)
		if ok {
			isOwnProfile = viewerID == user.ID
			if !isOwnProfile {
				isFollowing, _ = h.service.IsFollowing(
					c.Request.Context(),
					viewerID,
					user.ID,
				)
			}
		}
	}

	c.JSON(
		http.StatusOK,
		ToUserProfileResponse(
			*user,
			isFollowing,
			isOwnProfile,
		),
	)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	var req UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user id",
		})
		return
	}

	user, err := h.service.UpdateUser(
		c.Request.Context(),
		userID,
		req,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(
		http.StatusOK,
		ToUserResponse(*user),
	)
}
