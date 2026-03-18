package controllers

import (
	"net/http"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/apperrors"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/auth"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/httpapi"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService  services.UserService
	tokenManager auth.TokenManager
}

type registerRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type userResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewUserController(userService services.UserService, tokenManager auth.TokenManager) *UserController {
	return &UserController{userService: userService, tokenManager: tokenManager}
}

// Register handles user registration
// @Summary Register a new user
// @Description Create a user account and return a signed auth token.
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body registerRequest true "Registration payload"
// @Success 201 {object} authResponse
// @Failure 400 {object} httpapi.ErrorResponse
// @Failure 409 {object} httpapi.ErrorResponse
// @Failure 500 {object} httpapi.ErrorResponse
// @Router /api/v1/register [post]
func (uc *UserController) Register(c *gin.Context) {
	ctx := c.Request.Context()
	var req registerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		httpapi.WriteError(c, apperrors.Validation("invalid_request", "invalid request payload"))
		return
	}

	createdUser, err := uc.userService.RegisterUser(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		httpapi.WriteError(c, err)
		return
	}

	token, err := uc.tokenManager.GenerateToken(createdUser)
	if err != nil {
		httpapi.WriteError(c, apperrors.Internal("token_generation_failed", "failed to generate token", err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered",
		"token":   token,
		"user":    newUserResponse(createdUser),
	})
}

// Login handles user authentication
// @Summary Log in a user
// @Description Authenticate a user and return a signed auth token.
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body loginRequest true "Login payload"
// @Success 200 {object} authResponse
// @Failure 400 {object} httpapi.ErrorResponse
// @Failure 401 {object} httpapi.ErrorResponse
// @Failure 500 {object} httpapi.ErrorResponse
// @Router /api/v1/login [post]
func (uc *UserController) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		httpapi.WriteError(c, apperrors.Validation("invalid_request", "invalid request payload"))
		return
	}

	user, err := uc.userService.AuthenticateUser(ctx, req.Email, req.Password)
	if err != nil {
		httpapi.WriteError(c, err)
		return
	}

	token, err := uc.tokenManager.GenerateToken(user)
	if err != nil {
		httpapi.WriteError(c, apperrors.Internal("token_generation_failed", "failed to generate token", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user":    newUserResponse(user),
	})
}

func newUserResponse(user *models.User) userResponse {
	return userResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}
