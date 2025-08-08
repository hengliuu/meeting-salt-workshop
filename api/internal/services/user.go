package services

import (
	"api/internal/models"
	"api/internal/repositories"
	"api/internal/utils"
	"errors"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.New("validation failed")
	}

	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	existingUserByMSID, err := s.userRepo.GetByMicrosoftID(req.MicrosoftID)
	if err == nil && existingUserByMSID != nil {
		return nil, errors.New("user with this Microsoft ID already exists")
	}

	user := &models.User{
		MicrosoftID:    req.MicrosoftID,
		Email:          req.Email,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		DisplayName:    req.DisplayName,
		ProfilePicture: req.ProfilePicture,
		Role:           models.RoleEmployee,
		IsActive:       true,
	}

	return s.userRepo.Create(user)
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.GetByEmail(email)
}

func (s *UserService) GetUserByMicrosoftID(microsoftID string) (*models.User, error) {
	return s.userRepo.GetByMicrosoftID(microsoftID)
}

func (s *UserService) GetAllUsers(page, limit int) ([]*models.User, utils.PaginationMeta, error) {
	offset := utils.GetOffset(page, limit)
	users, total, err := s.userRepo.GetAll(offset, limit)
	if err != nil {
		return nil, utils.PaginationMeta{}, err
	}

	meta := utils.CreatePaginationMeta(page, limit, total)
	return users, meta, nil
}

func (s *UserService) UpdateUser(id uint, req *models.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.DisplayName != nil {
		user.DisplayName = *req.DisplayName
	}
	if req.ProfilePicture != nil {
		user.ProfilePicture = *req.ProfilePicture
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	return s.userRepo.Update(user)
}

func (s *UserService) DeleteUser(id uint) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	if !user.IsActive {
		return errors.New("user is already inactive")
	}

	user.IsActive = false
	_, err = s.userRepo.Update(user)
	return err
}

func (s *UserService) GetActiveUsers() ([]*models.User, error) {
	return s.userRepo.GetActiveUsers()
}

func (s *UserService) SearchUsers(query string, page, limit int) ([]*models.User, utils.PaginationMeta, error) {
	offset := utils.GetOffset(page, limit)
	users, total, err := s.userRepo.SearchUsers(query, offset, limit)
	if err != nil {
		return nil, utils.PaginationMeta{}, err
	}

	meta := utils.CreatePaginationMeta(page, limit, total)
	return users, meta, nil
}

func (s *UserService) UpdateUserRole(id uint, role models.UserRole) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	user.Role = role
	_, err = s.userRepo.Update(user)
	return err
}

func (s *UserService) ActivateUser(id uint) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	if user.IsActive {
		return errors.New("user is already active")
	}

	user.IsActive = true
	_, err = s.userRepo.Update(user)
	return err
}

func (s *UserService) DeactivateUser(id uint) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	if !user.IsActive {
		return errors.New("user is already inactive")
	}

	user.IsActive = false
	_, err = s.userRepo.Update(user)
	return err
}