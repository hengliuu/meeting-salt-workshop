package models

import (
	"time"

	"gorm.io/gorm"
)

type Room struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Capacity    int            `json:"capacity" gorm:"not null"`
	Location    string         `json:"location"`
	Features    []RoomFeature  `json:"features" gorm:"many2many:room_room_features;"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	Meetings []Meeting `json:"meetings" gorm:"foreignKey:RoomID"`
}

type RoomFeature struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateRoomRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Capacity    int    `json:"capacity" validate:"required,min=1"`
	Location    string `json:"location"`
	FeatureIDs  []uint `json:"feature_ids"`
}

type UpdateRoomRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Capacity    *int    `json:"capacity" validate:"omitempty,min=1"`
	Location    *string `json:"location"`
	FeatureIDs  []uint  `json:"feature_ids"`
	IsActive    *bool   `json:"is_active"`
}

type CreateRoomFeatureRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type UpdateRoomFeatureRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type RoomAvailabilityQuery struct {
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required"`
	Capacity  *int      `json:"capacity"`
}