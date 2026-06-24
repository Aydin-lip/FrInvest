package service

import (
	"errors"
	"fmt"
	"math/rand"
	"recruitment-api/internal/dto"
	"recruitment-api/internal/email"
	"recruitment-api/internal/models"
	"recruitment-api/internal/repository"
	"regexp"
	"time"

	"gorm.io/gorm"
)

var phoneRegex = regexp.MustCompile(`^09\d{9}$`)

type AuthService interface {
	SendCode(req dto.SendCodeRequest) error
	Register(req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
}

type authService struct {
	userRepo         repository.UserRepository
	verificationRepo repository.VerificationRepository
	jwtService       JWTService
	mailer           *email.Mailer
}

func NewAuthService(
	userRepo repository.UserRepository,
	verificationRepo repository.VerificationRepository,
	jwtService JWTService,
	mailer *email.Mailer,
) AuthService {
	return &authService{
		userRepo:         userRepo,
		verificationRepo: verificationRepo,
		jwtService:       jwtService,
		mailer:           mailer,
	}
}

func (s *authService) SendCode(req dto.SendCodeRequest) error {
	// Check for an active (non-expired, non-used) code
	existing, err := s.verificationRepo.FindActiveByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("database error")
	}
	if existing != nil {
		return errors.New("please wait before requesting another code")
	}

	// Generate 6-digit random code
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	vc := &models.VerificationCode{
		Email:     req.Email,
		Code:      code,
		ExpiresAt: time.Now().Add(2 * time.Minute),
		IsUsed:    false,
	}

	if err := s.verificationRepo.Create(vc); err != nil {
		return fmt.Errorf("failed to store verification code")
	}

	// Send email (non-blocking failure is acceptable — log in production)
	_ = s.mailer.SendVerificationCode(req.Email, code)

	return nil
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Validate phone number if provided
	if req.PhoneNumber != nil && *req.PhoneNumber != "" {
		if !phoneRegex.MatchString(*req.PhoneNumber) {
			return nil, errors.New("invalid phone number format")
		}
	}

	// Check user does not already exist
	existing, err := s.userRepo.FindByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("database error")
	}
	if existing != nil {
		return nil, errors.New("user already exists with this email")
	}

	// Validate verification code
	vc, err := s.verificationRepo.FindValidCode(req.Email, req.Code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid, expired, or already used verification code")
		}
		return nil, fmt.Errorf("database error")
	}

	// Create user
	user := &models.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Status:      0,
		IsActive:    true,
		IsDeleted:   false,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user")
	}

	// Mark code as used
	if err := s.verificationRepo.MarkAsUsed(vc.ID); err != nil {
		return nil, fmt.Errorf("failed to mark verification code as used")
	}

	// Generate JWT
	token, err := s.jwtService.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token")
	}

	// Send welcome email
	_ = s.mailer.SendWelcomeEmail(user.Email, user.FirstName)

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:          user.ID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			PhoneNumber: user.PhoneNumber,
			Email:       user.Email,
			Status:      user.Status,
		},
	}, nil
}

func (s *authService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error")
	}

	// Validate verification code
	vc, err := s.verificationRepo.FindValidCode(req.Email, req.Code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid, expired, or already used verification code")
		}
		return nil, fmt.Errorf("database error")
	}

	// Mark code as used
	if err := s.verificationRepo.MarkAsUsed(vc.ID); err != nil {
		return nil, fmt.Errorf("failed to mark verification code as used")
	}

	// Generate JWT
	token, err := s.jwtService.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token")
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Status:    user.Status,
		},
	}, nil
}
