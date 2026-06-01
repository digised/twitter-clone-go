package user

import (
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
		panic("user service is null")
	}

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

	if viewerID, ok := getUserID(c); ok {
		isOwnProfile = viewerID == user.ID

		if !isOwnProfile {
			isFollowing, err = h.service.IsFollowing(
				c.Request.Context(),
				viewerID,
				user.ID,
			)

			if err != nil {
				isFollowing = false
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

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
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

func getUserID(c *gin.Context) (uuid.UUID, bool) {
	value, exists := c.Get(ContextUserIDKey)
	if !exists {
		return uuid.Nil, false
	}

	userID, ok := value.(uuid.UUID)
	if !ok {
		return uuid.Nil, false
	}

	return userID, true
}
