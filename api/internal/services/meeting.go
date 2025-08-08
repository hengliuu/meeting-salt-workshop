package services

import (
	"api/internal/models"
	"api/internal/repositories"
	"api/internal/utils"
	"errors"
	"time"
)

type MeetingService struct {
	meetingRepo repositories.MeetingRepository
	roomRepo    repositories.RoomRepository
	userRepo    repositories.UserRepository
}

func NewMeetingService(meetingRepo repositories.MeetingRepository, roomRepo repositories.RoomRepository, userRepo repositories.UserRepository) *MeetingService {
	return &MeetingService{
		meetingRepo: meetingRepo,
		roomRepo:    roomRepo,
		userRepo:    userRepo,
	}
}

func (s *MeetingService) CreateMeeting(organizerID uint, req *models.CreateMeetingRequest) (*models.Meeting, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.New("validation failed")
	}

	if req.EndTime.Before(req.StartTime) {
		return nil, errors.New("end time must be after start time")
	}

	if req.StartTime.Before(time.Now()) {
		return nil, errors.New("meeting start time cannot be in the past")
	}

	room, err := s.roomRepo.GetByID(req.RoomID)
	if err != nil {
		return nil, errors.New("room not found")
	}

	if !room.IsActive {
		return nil, errors.New("room is not active")
	}

	available, err := s.roomRepo.IsRoomAvailable(req.RoomID, req.StartTime, req.EndTime, nil)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, errors.New("room is not available for the selected time")
	}

	organizer, err := s.userRepo.GetByID(organizerID)
	if err != nil {
		return nil, errors.New("organizer not found")
	}

	if !organizer.IsActive {
		return nil, errors.New("organizer is not active")
	}

	meeting := &models.Meeting{
		Title:             req.Title,
		Description:       req.Description,
		StartTime:         req.StartTime,
		EndTime:           req.EndTime,
		OrganizerID:       organizerID,
		RoomID:            req.RoomID,
		Status:            models.StatusScheduled,
		IsRecurring:       req.IsRecurring,
		RecurrencePattern: req.RecurrencePattern,
	}

	createdMeeting, err := s.meetingRepo.Create(meeting)
	if err != nil {
		return nil, err
	}

	if len(req.AttendeeIDs) > 0 {
		for _, attendeeID := range req.AttendeeIDs {
			if attendeeID == organizerID {
				continue
			}

			attendee, err := s.userRepo.GetByID(attendeeID)
			if err != nil || !attendee.IsActive {
				continue
			}

			if err := s.meetingRepo.AddAttendee(createdMeeting.ID, attendeeID); err != nil {
				continue
			}
		}
	}

	return s.meetingRepo.GetByID(createdMeeting.ID)
}

func (s *MeetingService) GetMeetingByID(id uint) (*models.Meeting, error) {
	return s.meetingRepo.GetByID(id)
}

func (s *MeetingService) GetAllMeetings(page, limit int) ([]*models.Meeting, utils.PaginationMeta, error) {
	offset := utils.GetOffset(page, limit)
	meetings, total, err := s.meetingRepo.GetAll(offset, limit)
	if err != nil {
		return nil, utils.PaginationMeta{}, err
	}

	meta := utils.CreatePaginationMeta(page, limit, total)
	return meetings, meta, nil
}

func (s *MeetingService) UpdateMeeting(id uint, userID uint, req *models.UpdateMeetingRequest) (*models.Meeting, error) {
	meeting, err := s.meetingRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("meeting not found")
	}

	if meeting.OrganizerID != userID {
		user, err := s.userRepo.GetByID(userID)
		if err != nil || !user.IsManager() {
			return nil, errors.New("only organizer or manager can update meeting")
		}
	}

	if meeting.Status == models.StatusCompleted || meeting.Status == models.StatusCancelled {
		return nil, errors.New("cannot update completed or cancelled meeting")
	}

	if req.Title != nil {
		meeting.Title = *req.Title
	}
	if req.Description != nil {
		meeting.Description = *req.Description
	}

	if req.StartTime != nil || req.EndTime != nil {
		startTime := meeting.StartTime
		endTime := meeting.EndTime

		if req.StartTime != nil {
			if req.StartTime.Before(time.Now()) {
				return nil, errors.New("meeting start time cannot be in the past")
			}
			startTime = *req.StartTime
		}
		if req.EndTime != nil {
			endTime = *req.EndTime
		}

		if endTime.Before(startTime) {
			return nil, errors.New("end time must be after start time")
		}

		roomID := meeting.RoomID
		if req.RoomID != nil {
			roomID = *req.RoomID
		}

		available, err := s.roomRepo.IsRoomAvailable(roomID, startTime, endTime, &meeting.ID)
		if err != nil {
			return nil, err
		}
		if !available {
			return nil, errors.New("room is not available for the selected time")
		}

		meeting.StartTime = startTime
		meeting.EndTime = endTime
	}

	if req.RoomID != nil {
		room, err := s.roomRepo.GetByID(*req.RoomID)
		if err != nil {
			return nil, errors.New("room not found")
		}
		if !room.IsActive {
			return nil, errors.New("room is not active")
		}
		meeting.RoomID = *req.RoomID
	}

	if req.Status != nil {
		meeting.Status = *req.Status
	}

	if req.IsRecurring != nil {
		meeting.IsRecurring = *req.IsRecurring
	}

	if req.RecurrencePattern != nil {
		meeting.RecurrencePattern = *req.RecurrencePattern
	}

	updatedMeeting, err := s.meetingRepo.Update(meeting)
	if err != nil {
		return nil, err
	}

	if req.AttendeeIDs != nil {
		currentAttendees, _ := s.meetingRepo.GetMeetingAttendees(meeting.ID)
		currentAttendeeMap := make(map[uint]bool)
		for _, attendee := range currentAttendees {
			currentAttendeeMap[attendee.ID] = true
		}

		newAttendeeMap := make(map[uint]bool)
		for _, attendeeID := range req.AttendeeIDs {
			if attendeeID != meeting.OrganizerID {
				newAttendeeMap[attendeeID] = true
			}
		}

		for attendeeID := range currentAttendeeMap {
			if !newAttendeeMap[attendeeID] {
				s.meetingRepo.RemoveAttendee(meeting.ID, attendeeID)
			}
		}

		for attendeeID := range newAttendeeMap {
			if !currentAttendeeMap[attendeeID] {
				attendee, err := s.userRepo.GetByID(attendeeID)
				if err == nil && attendee.IsActive {
					s.meetingRepo.AddAttendee(meeting.ID, attendeeID)
				}
			}
		}
	}

	return s.meetingRepo.GetByID(updatedMeeting.ID)
}

func (s *MeetingService) DeleteMeeting(id uint, userID uint) error {
	meeting, err := s.meetingRepo.GetByID(id)
	if err != nil {
		return errors.New("meeting not found")
	}

	if meeting.OrganizerID != userID {
		user, err := s.userRepo.GetByID(userID)
		if err != nil || !user.IsManager() {
			return errors.New("only organizer or manager can delete meeting")
		}
	}

	if meeting.Status == models.StatusCompleted {
		return errors.New("cannot delete completed meeting")
	}

	return s.meetingRepo.UpdateMeetingStatus(id, models.StatusCancelled)
}

func (s *MeetingService) GetMeetingsByFilter(filter models.MeetingFilter, page, limit int) ([]*models.Meeting, utils.PaginationMeta, error) {
	offset := utils.GetOffset(page, limit)
	meetings, total, err := s.meetingRepo.GetByFilter(filter, offset, limit)
	if err != nil {
		return nil, utils.PaginationMeta{}, err
	}

	meta := utils.CreatePaginationMeta(page, limit, total)
	return meetings, meta, nil
}

func (s *MeetingService) GetUpcomingMeetings(userID *uint, limit int) ([]*models.Meeting, error) {
	return s.meetingRepo.GetUpcomingMeetings(userID, limit)
}

func (s *MeetingService) GetMeetingsByDateRange(startDate, endDate time.Time, userID *uint) ([]*models.Meeting, error) {
	if endDate.Before(startDate) {
		return nil, errors.New("end date must be after start date")
	}

	return s.meetingRepo.GetMeetingsByDateRange(startDate, endDate, userID)
}