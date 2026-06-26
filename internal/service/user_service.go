package service

import (
	"errors"
	"fmt"
	"math"
	"recruitment-api/internal/dto"
	"recruitment-api/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrStatusOutOfRange  = errors.New("status must be between 0 and 4")
	ErrUserNotVerified   = errors.New("user is not verified")
)

type UserService interface {
	GetVerifiedUsers() ([]dto.UserListItem, error)
	UpdateStatus(req dto.UpdateStatusRequest) error
	GetStatusPercentages() (*dto.StatusPercentagesResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetVerifiedUsers() ([]dto.UserListItem, error) {
	users, err := s.userRepo.GetVerifiedUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users")
	}

	result := make([]dto.UserListItem, 0, len(users))
	for _, u := range users {
		result = append(result, dto.UserListItem{
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Status:    u.Status,
		})
	}

	return result, nil
}

func (s *userService) UpdateStatus(req dto.UpdateStatusRequest) error {
	if req.Status < 0 || req.Status > 4 {
		return ErrStatusOutOfRange
	}

	user, err := s.userRepo.FindByID(req.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("database error")
	}

	if !user.Verify {
		return ErrUserNotVerified
	}

	if err := s.userRepo.UpdateStatus(req.UserID, req.Status); err != nil {
		return fmt.Errorf("failed to update status")
	}

	return nil
}

func (s *userService) GetStatusPercentages() (*dto.StatusPercentagesResponse, error) {
	total, err := s.userRepo.GetVerifiedTotalCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get user count")
	}

	if total == 0 {
		return &dto.StatusPercentagesResponse{}, nil
	}

	counts, err := s.userRepo.GetVerifiedStatusCounts()
	if err != nil {
		return nil, fmt.Errorf("failed to get status counts")
	}

	calc := func(status int8) float64 {
		return math.Round((float64(counts[status]) / float64(total)) * 100)
	}

	return &dto.StatusPercentagesResponse{
		New:         calc(0),
		Reviewed:    calc(1),
		Interviewed: calc(2),
		OfferSent:   calc(3),
		Rejected:    calc(4),
	}, nil
}
