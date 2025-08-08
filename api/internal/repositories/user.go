package repositories

import (
	"api/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByMicrosoftID(microsoftID string) (*models.User, error)
	GetAll(offset, limit int) ([]*models.User, int64, error)
	Update(user *models.User) (*models.User, error)
	Delete(id uint) error
	GetActiveUsers() ([]*models.User, error)
	SearchUsers(query string, offset, limit int) ([]*models.User, int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("OrganizedMeetings").Preload("AttendedMeetings").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByMicrosoftID(microsoftID string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("microsoft_id = ?", microsoftID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetAll(offset, limit int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	if err := r.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) Update(user *models.User) (*models.User, error) {
	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *userRepository) GetActiveUsers() ([]*models.User, error) {
	var users []*models.User
	if err := r.db.Where("is_active = ?", true).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) SearchUsers(query string, offset, limit int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	searchQuery := "%" + query + "%"
	dbQuery := r.db.Model(&models.User{}).Where(
		"first_name LIKE ? OR last_name LIKE ? OR email LIKE ? OR display_name LIKE ?",
		searchQuery, searchQuery, searchQuery, searchQuery,
	)

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := dbQuery.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}