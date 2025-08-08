package repositories

import (
	"api/internal/models"
	"time"

	"gorm.io/gorm"
)

type MeetingRepository interface {
	Create(meeting *models.Meeting) (*models.Meeting, error)
	GetByID(id uint) (*models.Meeting, error)
	GetAll(offset, limit int) ([]*models.Meeting, int64, error)
	Update(meeting *models.Meeting) (*models.Meeting, error)
	Delete(id uint) error
	GetByFilter(filter models.MeetingFilter, offset, limit int) ([]*models.Meeting, int64, error)
	GetUpcomingMeetings(userID *uint, limit int) ([]*models.Meeting, error)
	GetMeetingsByDateRange(startDate, endDate time.Time, userID *uint) ([]*models.Meeting, error)
	GetMeetingsByRoom(roomID uint, offset, limit int) ([]*models.Meeting, int64, error)
	GetMeetingsByOrganizer(organizerID uint, offset, limit int) ([]*models.Meeting, int64, error)
	GetConflictingMeetings(roomID uint, startTime, endTime time.Time, excludeMeetingID *uint) ([]*models.Meeting, error)
	UpdateMeetingStatus(id uint, status models.MeetingStatus) error
	GetMeetingAttendees(meetingID uint) ([]*models.User, error)
	AddAttendee(meetingID, userID uint) error
	RemoveAttendee(meetingID, userID uint) error
}

type meetingRepository struct {
	db *gorm.DB
}

func NewMeetingRepository(db *gorm.DB) MeetingRepository {
	return &meetingRepository{db: db}
}

func (r *meetingRepository) Create(meeting *models.Meeting) (*models.Meeting, error) {
	if err := r.db.Create(meeting).Error; err != nil {
		return nil, err
	}
	return r.GetByID(meeting.ID)
}

func (r *meetingRepository) GetByID(id uint) (*models.Meeting, error) {
	var meeting models.Meeting
	if err := r.db.Preload("Organizer").Preload("Room").Preload("Room.Features").Preload("Attendees").First(&meeting, id).Error; err != nil {
		return nil, err
	}
	return &meeting, nil
}

func (r *meetingRepository) GetAll(offset, limit int) ([]*models.Meeting, int64, error) {
	var meetings []*models.Meeting
	var total int64

	if err := r.db.Model(&models.Meeting{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("Organizer").Preload("Room").Preload("Attendees").
		Offset(offset).Limit(limit).Order("start_time DESC").Find(&meetings).Error; err != nil {
		return nil, 0, err
	}

	return meetings, total, nil
}

func (r *meetingRepository) Update(meeting *models.Meeting) (*models.Meeting, error) {
	if err := r.db.Save(meeting).Error; err != nil {
		return nil, err
	}
	return r.GetByID(meeting.ID)
}

func (r *meetingRepository) Delete(id uint) error {
	return r.db.Delete(&models.Meeting{}, id).Error
}

func (r *meetingRepository) GetByFilter(filter models.MeetingFilter, offset, limit int) ([]*models.Meeting, int64, error) {
	query := r.db.Model(&models.Meeting{})

	if filter.OrganizerID != nil {
		query = query.Where("organizer_id = ?", *filter.OrganizerID)
	}

	if filter.RoomID != nil {
		query = query.Where("room_id = ?", *filter.RoomID)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.StartDate != nil {
		query = query.Where("start_time >= ?", *filter.StartDate)
	}

	if filter.EndDate != nil {
		query = query.Where("end_time <= ?", *filter.EndDate)
	}

	if filter.UserID != nil {
		query = query.Where("organizer_id = ? OR id IN (SELECT meeting_id FROM meeting_attendees WHERE user_id = ?)", 
			*filter.UserID, *filter.UserID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var meetings []*models.Meeting
	if err := query.Preload("Organizer").Preload("Room").Preload("Attendees").
		Offset(offset).Limit(limit).Order("start_time DESC").Find(&meetings).Error; err != nil {
		return nil, 0, err
	}

	return meetings, total, nil
}

func (r *meetingRepository) GetUpcomingMeetings(userID *uint, limit int) ([]*models.Meeting, error) {
	query := r.db.Where("start_time > ? AND status = ?", time.Now(), models.StatusScheduled)

	if userID != nil {
		query = query.Where("organizer_id = ? OR id IN (SELECT meeting_id FROM meeting_attendees WHERE user_id = ?)", 
			*userID, *userID)
	}

	var meetings []*models.Meeting
	if err := query.Preload("Organizer").Preload("Room").Preload("Attendees").
		Order("start_time ASC").Limit(limit).Find(&meetings).Error; err != nil {
		return nil, err
	}

	return meetings, nil
}

func (r *meetingRepository) GetMeetingsByDateRange(startDate, endDate time.Time, userID *uint) ([]*models.Meeting, error) {
	query := r.db.Where("start_time >= ? AND end_time <= ?", startDate, endDate)

	if userID != nil {
		query = query.Where("organizer_id = ? OR id IN (SELECT meeting_id FROM meeting_attendees WHERE user_id = ?)", 
			*userID, *userID)
	}

	var meetings []*models.Meeting
	if err := query.Preload("Organizer").Preload("Room").Preload("Attendees").
		Order("start_time ASC").Find(&meetings).Error; err != nil {
		return nil, err
	}

	return meetings, nil
}

func (r *meetingRepository) GetMeetingsByRoom(roomID uint, offset, limit int) ([]*models.Meeting, int64, error) {
	var meetings []*models.Meeting
	var total int64

	query := r.db.Model(&models.Meeting{}).Where("room_id = ?", roomID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Organizer").Preload("Room").Preload("Attendees").
		Offset(offset).Limit(limit).Order("start_time DESC").Find(&meetings).Error; err != nil {
		return nil, 0, err
	}

	return meetings, total, nil
}

func (r *meetingRepository) GetMeetingsByOrganizer(organizerID uint, offset, limit int) ([]*models.Meeting, int64, error) {
	var meetings []*models.Meeting
	var total int64

	query := r.db.Model(&models.Meeting{}).Where("organizer_id = ?", organizerID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Organizer").Preload("Room").Preload("Attendees").
		Offset(offset).Limit(limit).Order("start_time DESC").Find(&meetings).Error; err != nil {
		return nil, 0, err
	}

	return meetings, total, nil
}

func (r *meetingRepository) GetConflictingMeetings(roomID uint, startTime, endTime time.Time, excludeMeetingID *uint) ([]*models.Meeting, error) {
	query := r.db.Where("room_id = ? AND start_time < ? AND end_time > ? AND status IN ?", 
		roomID, endTime, startTime, []models.MeetingStatus{models.StatusScheduled, models.StatusInProgress})

	if excludeMeetingID != nil {
		query = query.Where("id != ?", *excludeMeetingID)
	}

	var meetings []*models.Meeting
	if err := query.Preload("Organizer").Preload("Room").Find(&meetings).Error; err != nil {
		return nil, err
	}

	return meetings, nil
}

func (r *meetingRepository) UpdateMeetingStatus(id uint, status models.MeetingStatus) error {
	return r.db.Model(&models.Meeting{}).Where("id = ?", id).Update("status", status).Error
}

func (r *meetingRepository) GetMeetingAttendees(meetingID uint) ([]*models.User, error) {
	var meeting models.Meeting
	if err := r.db.Preload("Attendees").First(&meeting, meetingID).Error; err != nil {
		return nil, err
	}

	users := make([]*models.User, len(meeting.Attendees))
	for i := range meeting.Attendees {
		users[i] = &meeting.Attendees[i]
	}

	return users, nil
}

func (r *meetingRepository) AddAttendee(meetingID, userID uint) error {
	return r.db.Exec("INSERT IGNORE INTO meeting_attendees (meeting_id, user_id) VALUES (?, ?)", 
		meetingID, userID).Error
}

func (r *meetingRepository) RemoveAttendee(meetingID, userID uint) error {
	return r.db.Exec("DELETE FROM meeting_attendees WHERE meeting_id = ? AND user_id = ?", 
		meetingID, userID).Error
}