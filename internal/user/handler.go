package user

import (
	"errors"
	"net/http"

	"locallyn-be/internal/common/auth"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SignUp(c *gin.Context) {
	var req SignUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.service.SignUp(c.Request.Context(), req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SignUpResponse{
		Message: "User Created Successfully",
	})
}
func (h *Handler) VerifyUser(c *gin.Context) {
	var req VerifyUserRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.service.VerifyUser(c.Request.Context(), req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, verifyUserResponse{
		Message: "User Verified Successfully",
	})

}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := h.service.Login(c.Request.Context(), req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, LoginResponse{
		AccessToken: token,
	})

}

func (h *Handler) CreateProfile(c *gin.Context) {
	var req CreateProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	claims, err := auth.GetClaimsFromContext(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err = h.service.CreateProfile(c.Request.Context(), claims, req)
	if err != nil {
		if errors.Is(err, ErrUsernameAlreadyExists) || errors.Is(err, ErrUserProfileExists) {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, CreateProfileResponse{
		Message: "Profile Created Successfully",
	})
}
