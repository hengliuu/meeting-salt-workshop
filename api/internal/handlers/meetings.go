package handlers

import (
	"api/internal/models"
	"api/internal/services"
	"api/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MeetingHandler struct {
	meetingService *services.MeetingService
}

func NewMeetingHandler(meetingService *services.MeetingService) *MeetingHandler {
	return &MeetingHandler{
		meetingService: meetingService,
	}
}

func (h *MeetingHandler) CreateMeeting(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}
	
	organizerID, ok := userID.(uint)
	if !ok {
		utils.UnauthorizedResponse(c, "Invalid user ID")
		return
	}
	
	var req models.CreateMeetingRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request format")
		return
	}
	
	if validationErrors := utils.ValidateStruct(&req); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}
	
	meeting, err := h.meetingService.CreateMeeting(organizerID, &req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.CreatedResponse(c, "Meeting created successfully", meeting)
}

func (h *MeetingHandler) GetMeeting(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid meeting ID")
		return
	}
	
	meeting, err := h.meetingService.GetMeetingByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Meeting not found")
		return
	}
	
	utils.SuccessResponse(c, "Meeting retrieved successfully", meeting)
}

func (h *MeetingHandler) GetAllMeetings(c *gin.Context) {
	pagination := utils.GetPaginationParams(c)
	
	meetings, meta, err := h.meetingService.GetAllMeetings(pagination.Page, pagination.Limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve meetings")
		return
	}
	
	utils.PaginatedSuccessResponse(c, "Meetings retrieved successfully", meetings, meta)
}

func (h *MeetingHandler) UpdateMeeting(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid meeting ID")
		return
	}
	
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}
	
	currentUserID, ok := userID.(uint)
	if !ok {
		utils.UnauthorizedResponse(c, "Invalid user ID")
		return
	}
	
	var req models.UpdateMeetingRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request format")
		return
	}
	
	meeting, err := h.meetingService.UpdateMeeting(uint(id), currentUserID, &req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, "Meeting updated successfully", meeting)
}

func (h *MeetingHandler) DeleteMeeting(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid meeting ID")
		return
	}
	
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}
	
	currentUserID, ok := userID.(uint)
	if !ok {
		utils.UnauthorizedResponse(c, "Invalid user ID")
		return
	}
	
	if err := h.meetingService.DeleteMeeting(uint(id), currentUserID); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, "Meeting deleted successfully", nil)
}

func (h *MeetingHandler) GetMeetingsByFilter(c *gin.Context) {
	var filter models.MeetingFilter
	
	if organizerIDStr := c.Query("organizer_id"); organizerIDStr != "" {
		organizerID, err := strconv.ParseUint(organizerIDStr, 10, 32)
		if err != nil {
			utils.BadRequestResponse(c, "Invalid organizer_id")
			return
		}
		organizerIDUint := uint(organizerID)
		filter.OrganizerID = &organizerIDUint
	}
	
	if roomIDStr := c.Query("room_id"); roomIDStr != "" {
		roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
		if err != nil {
			utils.BadRequestResponse(c, "Invalid room_id")
			return
		}
		roomIDUint := uint(roomID)
		filter.RoomID = &roomIDUint
	}
	
	if statusStr := c.Query("status"); statusStr != "" {
		status := models.MeetingStatus(statusStr)
		filter.Status = &status
	}
	
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			utils.BadRequestResponse(c, "Invalid start_date format. Use YYYY-MM-DD")
			return
		}
		filter.StartDate = &startDate
	}
	
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			utils.BadRequestResponse(c, "Invalid end_date format. Use YYYY-MM-DD")
			return
		}
		filter.EndDate = &endDate
	}
	
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			utils.BadRequestResponse(c, "Invalid user_id")
			return
		}
		userIDUint := uint(userID)
		filter.UserID = &userIDUint
	}
	
	pagination := utils.GetPaginationParams(c)
	
	meetings, meta, err := h.meetingService.GetMeetingsByFilter(filter, pagination.Page, pagination.Limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve meetings")
		return
	}
	
	utils.PaginatedSuccessResponse(c, "Filtered meetings retrieved successfully", meetings, meta)
}

func (h *MeetingHandler) GetUpcomingMeetings(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	
	var userID *uint
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		userIDVal, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			utils.BadRequestResponse(c, "Invalid user_id")
			return
		}
		userIDUint := uint(userIDVal)
		userID = &userIDUint
	}
	
	meetings, err := h.meetingService.GetUpcomingMeetings(userID, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve upcoming meetings")
		return
	}
	
	utils.SuccessResponse(c, "Upcoming meetings retrieved successfully", meetings)
}

func (h *MeetingHandler) GetMeetingsByDateRange(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	
	if startDateStr == "" || endDateStr == "" {
		utils.BadRequestResponse(c, "start_date and end_date are required")
		return
	}
	
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid start_date format. Use YYYY-MM-DD")
		return
	}
	
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid end_date format. Use YYYY-MM-DD")
		return
	}
	
	var userID *uint
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		userIDVal, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			utils.BadRequestResponse(c, "Invalid user_id")
			return
		}
		userIDUint := uint(userIDVal)
		userID = &userIDUint
	}
	
	meetings, err := h.meetingService.GetMeetingsByDateRange(startDate, endDate, userID)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, "Meetings by date range retrieved successfully", meetings)
}

func (h *MeetingHandler) StartMeeting(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid meeting ID")
		return
	}
	
	if err := h.meetingService.StartMeeting(uint(id)); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, "Meeting started successfully", nil)
}

func (h *MeetingHandler) CompleteMeeting(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid meeting ID")
		return
	}
	
	if err := h.meetingService.CompleteMeeting(uint(id)); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, "Meeting completed successfully", nil)
}

func (h *MeetingHandler) CancelMeeting(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid meeting ID")
		return
	}
	
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}
	
	currentUserID, ok := userID.(uint)
	if !ok {
		utils.UnauthorizedResponse(c, "Invalid user ID")
		return
	}
	
	if err := h.meetingService.CancelMeeting(uint(id), currentUserID); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, "Meeting cancelled successfully", nil)
}

func (h *MeetingHandler) AddAttendee(c *gin.Context) {
	idParam := c.Param("id")
	meetingID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid meeting ID")
		return
	}
	
	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request format")
		return
	}
	
	if err := h.meetingService.AddAttendee(uint(meetingID), req.UserID); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, "Attendee added successfully", nil)
}

func (h *MeetingHandler) RemoveAttendee(c *gin.Context) {
	idParam := c.Param("id")
	meetingID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid meeting ID")
		return
	}
	
	userIDParam := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}
	
	if err := h.meetingService.RemoveAttendee(uint(meetingID), uint(userID)); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, "Attendee removed successfully", nil)
}

func (h *MeetingHandler) GetMeetingAttendees(c *gin.Context) {
	idParam := c.Param("id")
	meetingID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid meeting ID")
		return
	}
	
	attendees, err := h.meetingService.GetMeetingAttendees(uint(meetingID))
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, "Meeting attendees retrieved successfully", attendees)
}