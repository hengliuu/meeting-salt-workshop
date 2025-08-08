package models

import (
	"time"

	"gorm.io/gorm"
)

type Meeting struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	Title             string         `json:"title" gorm:"not null"`
	Description       string         `json:"description"`
	StartTime         time.Time      `json:"start_time" gorm:"not null"`
	EndTime           time.Time      `json:"end_time" gorm:"not null"`
	Status            MeetingStatus  `json:"status" gorm:"default:'scheduled'"`
	IsRecurring       bool           `json:"is_recurring" gorm:"default:false"`
	RecurrencePattern string         `json:"recurrence_pattern"`
	OrganizerID       uint           `json:"organizer_id" gorm:"not null"`
	RoomID            uint           `json:"room_id" gorm:"not null"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`

	Organizer User   `json:"organizer" gorm:"foreignKey:OrganizerID"`
	Room      Room   `json:"room" gorm:"foreignKey:RoomID"`
	Attendees []User `json:"attendees" gorm:"many2many:meeting_attendees;"`
}

type MeetingStatus string

const (
	StatusScheduled  MeetingStatus = "scheduled"
	StatusInProgress MeetingStatus = "in_progress"
	StatusCompleted  MeetingStatus = "completed"
	StatusCancelled  MeetingStatus = "cancelled"
)

type CreateMeetingRequest struct {
	Title             string    `json:"title" validate:"required"`
	Description       string    `json:"description"`
	StartTime         time.Time `json:"start_time" validate:"required"`
	EndTime           time.Time `json:"end_time" validate:"required"`
	RoomID            uint      `json:"room_id" validate:"required"`
	AttendeeIDs       []uint    `json:"attendee_ids"`
	IsRecurring       bool      `json:"is_recurring"`
	RecurrencePattern string    `json:"recurrence_pattern"`
}

type UpdateMeetingRequest struct {
	Title             *string        `json:"title"`
	Description       *string        `json:"description"`
	StartTime         *time.Time     `json:"start_time"`
	EndTime           *time.Time     `json:"end_time"`
	RoomID            *uint          `json:"room_id"`
	AttendeeIDs       []uint         `json:"attendee_ids"`
	Status            *MeetingStatus `json:"status"`
	IsRecurring       *bool          `json:"is_recurring"`
	RecurrencePattern *string        `json:"recurrence_pattern"`
}

type MeetingFilter struct {
	OrganizerID *uint          `json:"organizer_id"`
	RoomID      *uint          `json:"room_id"`
	Status      *MeetingStatus `json:"status"`
	StartDate   *time.Time     `json:"start_date"`
	EndDate     *time.Time     `json:"end_date"`
	UserID      *uint          `json:"user_id"`
}

func (m *Meeting) IsOverlapping(startTime, endTime time.Time) bool {
	return m.StartTime.Before(endTime) && m.EndTime.After(startTime)
}

func (m *Meeting) Duration() time.Duration {
	return m.EndTime.Sub(m.StartTime)
}

func (m *Meeting) IsActive() bool {
	return m.Status == StatusScheduled || m.Status == StatusInProgress
}