package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	MicrosoftID    string         `json:"microsoft_id" gorm:"uniqueIndex;not null"`
	Email          string         `json:"email" gorm:"uniqueIndex;not null"`
	FirstName      string         `json:"first_name" gorm:"not null"`
	LastName       string         `json:"last_name" gorm:"not null"`
	DisplayName    string         `json:"display_name"`
	ProfilePicture string         `json:"profile_picture"`
	Role           UserRole       `json:"role" gorm:"default:'employee'"`
	IsActive       bool           `json:"is_active" gorm:"default:true"`
	LastLogin      *time.Time     `json:"last_login"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`

	OrganizedMeetings []Meeting `json:"organized_meetings" gorm:"foreignKey:OrganizerID"`
	AttendedMeetings  []Meeting `json:"attended_meetings" gorm:"many2many:meeting_attendees;"`
}

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleManager  UserRole = "manager"
	RoleEmployee UserRole = "employee"
)

type CreateUserRequest struct {
	MicrosoftID    string `json:"microsoft_id" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	FirstName      string `json:"first_name" validate:"required"`
	LastName       string `json:"last_name" validate:"required"`
	DisplayName    string `json:"display_name"`
	ProfilePicture string `json:"profile_picture"`
}

type UpdateUserRequest struct {
	FirstName      *string   `json:"first_name"`
	LastName       *string   `json:"last_name"`
	DisplayName    *string   `json:"display_name"`
	ProfilePicture *string   `json:"profile_picture"`
	Role           *UserRole `json:"role"`
	IsActive       *bool     `json:"is_active"`
}

func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) IsManager() bool {
	return u.Role == RoleManager || u.Role == RoleAdmin
}