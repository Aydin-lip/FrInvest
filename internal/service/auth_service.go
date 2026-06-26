package service

import (
	"errors"
	"fmt"
	"recruitment-api/internal/dto"
	"recruitment-api/internal/email"
	"recruitment-api/internal/models"
	"recruitment-api/internal/repository"
	"recruitment-api/internal/utils"
	"regexp"
	"time"

	"gorm.io/gorm"
)

var phoneRegex = regexp.MustCompile(`^09\d{9}$`)

var (
	ErrUserAlreadyVerified = errors.New("User already exists and verified")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidPhoneNumber  = errors.New("invalid phone number format")
)

type RateLimitError struct {
	RemainingSeconds int
}

func (e *RateLimitError) Error() string {
	return "Please wait before requesting another email"
}

type VerifyEmailOutcome int

const (
	VerifySuccess VerifyEmailOutcome = iota
	VerifyExpired
	VerifyInvalid
)

type AuthService interface {
	Register(req dto.RegisterRequest) (string, error)
	SendVerification(req dto.SendVerificationRequest) error
	VerifyEmail(token string) (VerifyEmailOutcome, error)
}

type authService struct {
	userRepo         repository.UserRepository
	verificationRepo repository.VerificationRepository
	mailer           *email.Mailer
}

func NewAuthService(
	userRepo repository.UserRepository,
	verificationRepo repository.VerificationRepository,
	mailer *email.Mailer,
) AuthService {
	return &authService{
		userRepo:         userRepo,
		verificationRepo: verificationRepo,
		mailer:           mailer,
	}
}

func (s *authService) Register(req dto.RegisterRequest) (string, error) {
	if req.PhoneNumber != nil && *req.PhoneNumber != "" {
		if !phoneRegex.MatchString(*req.PhoneNumber) {
			return "", ErrInvalidPhoneNumber
		}
	}

	existing, err := s.userRepo.FindByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", fmt.Errorf("database error")
	}

	if existing != nil {
		if existing.Verify {
			return "", ErrUserAlreadyVerified
		}

		existing.FirstName = req.FirstName
		existing.LastName = req.LastName
		existing.PhoneNumber = req.PhoneNumber

		if err := s.userRepo.Update(existing); err != nil {
			return "", fmt.Errorf("failed to update user")
		}

		return s.issueRegistrationVerificationToken(existing.ID, existing.Email)
	}

	user := &models.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Status:      0,
		Verify:      false,
		IsActive:    true,
		IsDeleted:   false,
	}

	if err := s.userRepo.Create(user); err != nil {
		return "", fmt.Errorf("failed to create user")
	}

	return s.issueRegistrationVerificationToken(user.ID, user.Email)
}

func (s *authService) SendVerification(req dto.SendVerificationRequest) error {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("database error")
	}

	active, err := s.verificationRepo.FindLatestValidByUserID(user.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("database error")
	}
	if active != nil {
		remaining := int(time.Until(active.ExpiresAt).Seconds())
		if remaining < 0 {
			remaining = 0
		}
		return &RateLimitError{RemainingSeconds: remaining}
	}

	token, err := s.createVerificationToken(user.ID)
	if err != nil {
		return err
	}

	_ = s.mailer.SendVerificationEmail(user.Email, token)
	return nil
}

func (s *authService) VerifyEmail(token string) (VerifyEmailOutcome, error) {
	vt, err := s.verificationRepo.FindByToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return VerifyInvalid, nil
		}
		return VerifyInvalid, fmt.Errorf("database error")
	}

	if vt.UsedAt != nil {
		return VerifyInvalid, nil
	}

	if time.Now().After(vt.ExpiresAt) {
		return VerifyExpired, nil
	}

	latest, err := s.verificationRepo.FindLatestValidByUserID(vt.UserID)
	if err != nil || latest.ID != vt.ID {
		return VerifyInvalid, nil
	}

	now := time.Now()
	if err := s.verificationRepo.MarkAsUsed(vt.ID, now); err != nil {
		return VerifyInvalid, fmt.Errorf("failed to mark token as used")
	}

	if err := s.userRepo.SetVerified(vt.UserID); err != nil {
		return VerifyInvalid, fmt.Errorf("failed to verify user")
	}

	user, err := s.userRepo.FindByID(vt.UserID)
	if err == nil {
		_ = s.mailer.SendWebinarEmail(user.Email, user.FirstName)
	}

	return VerifySuccess, nil
}

func (s *authService) issueRegistrationVerificationToken(userID uint64, email string) (string, error) {
	token, err := s.resolveVerificationTokenForRegistration(userID)
	if err != nil {
		return "", err
	}

	_ = s.mailer.SendVerificationEmail(email, token)
	return token, nil
}

func (s *authService) resolveVerificationTokenForRegistration(userID uint64) (string, error) {
	existing, err := s.verificationRepo.FindLatestUnusedByUserID(userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", fmt.Errorf("database error")
	}

	if existing == nil {
		return s.createVerificationToken(userID)
	}

	if time.Now().Before(existing.ExpiresAt) {
		return existing.Token, nil
	}

	token, err := utils.GenerateToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate token")
	}

	expiresAt := time.Now().Add(2 * time.Minute)
	if err := s.verificationRepo.UpdateToken(existing.ID, token, expiresAt); err != nil {
		return "", fmt.Errorf("failed to update verification token")
	}

	return token, nil
}

func (s *authService) createVerificationToken(userID uint64) (string, error) {
	token, err := utils.GenerateToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate token")
	}

	vt := &models.VerificationToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(2 * time.Minute),
	}

	if err := s.verificationRepo.Create(vt); err != nil {
		return "", fmt.Errorf("failed to store verification token")
	}

	return token, nil
}
