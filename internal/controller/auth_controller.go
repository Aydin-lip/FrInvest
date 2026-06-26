package controller

import (
	"errors"
	"net/http"
	"recruitment-api/config"
	"recruitment-api/internal/dto"
	"recruitment-api/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ac *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	token, err := ac.authService.Register(req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserAlreadyVerified):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Message: err.Error()})
		case errors.Is(err, service.ErrInvalidPhoneNumber):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, dto.TokenResponse{Token: token})
}

func (ac *AuthController) SendVerification(c *gin.Context) {
	var req dto.SendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	err := ac.authService.SendVerification(req)
	if err != nil {
		var rateLimit *service.RateLimitError
		if errors.As(err, &rateLimit) {
			c.JSON(http.StatusTooManyRequests, dto.RateLimitResponse{
				Message:          rateLimit.Error(),
				RemainingSeconds: rateLimit.RemainingSeconds,
			})
			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Verification email sent successfully"})
}

func (ac *AuthController) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "token is required"})
		return
	}

	outcome, err := ac.authService.VerifyEmail(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
		return
	}

	switch outcome {
	case service.VerifySuccess:
		c.Redirect(http.StatusFound, config.AppConfig.FrontendSuccessURL)
	case service.VerifyExpired:
		c.Redirect(http.StatusFound, config.AppConfig.FrontendErrorURL)
	default:
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid or expired verification token"})
	}
}
