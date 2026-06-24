package service

import (
	"errors"
	"fmt"
	"recruitment-api/internal/dto"
	"recruitment-api/internal/repository"
)

type UserService interface {
	GetAll() ([]dto.UserListItem, error)
	UpdateStatus(req dto.UpdateStatusRequest) error
	GetStatistics() (*dto.StatisticsResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetAll() ([]dto.UserListItem, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users")
	}

	var result []dto.UserListItem
	for _, u := range users {
		result = append(result, dto.UserListItem{
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Status:    u.Status,
		})
	}

	// Return empty slice instead of nil
	if result == nil {
		result = []dto.UserListItem{}
	}

	return result, nil
}

func (s *userService) UpdateStatus(req dto.UpdateStatusRequest) error {
	if req.Status < 0 || req.Status > 4 {
		return errors.New("status must be between 0 and 4")
	}

	// Check user exists
	_, err := s.userRepo.FindByID(req.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	return s.userRepo.UpdateStatus(req.UserID, req.Status)
}

func (s *userService) GetStatistics() (*dto.StatisticsResponse, error) {
	total, err := s.userRepo.GetTotalCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get user count")
	}

	if total == 0 {
		return &dto.StatisticsResponse{
			New:         0,
			Reviewed:    0,
			Interviewed: 0,
			OfferSent:   0,
			Rejected:    0,
		}, nil
	}

	counts, err := s.userRepo.GetStatusCounts()
	if err != nil {
		return nil, fmt.Errorf("failed to get status counts")
	}

	calc := func(statusCount int64) float64 {
		return (float64(statusCount) * 100) / float64(total)
	}

	return &dto.StatisticsResponse{
		New:         calc(counts[0]),
		Reviewed:    calc(counts[1]),
		Interviewed: calc(counts[2]),
		OfferSent:   calc(counts[3]),
		Rejected:    calc(counts[4]),
	}, nil
}
