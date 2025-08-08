package handlers

import (
	"api/internal/models"
	"api/internal/services"
	"api/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardService *services.DashboardService
}

func NewDashboardHandler(dashboardService *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	filter := h.parseFilter(c)
	
	stats, err := h.dashboardService.GetDashboardStats(filter)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve dashboard stats")
		return
	}
	
	utils.SuccessResponse(c, "Dashboard stats retrieved successfully", stats)
}

func (h *DashboardHandler) GetRoomUtilization(c *gin.Context) {
	filter := h.parseFilter(c)
	
	utilization, err := h.dashboardService.GetRoomUtilization(filter)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve room utilization")
		return
	}
	
	utils.SuccessResponse(c, "Room utilization retrieved successfully", utilization)
}

func (h *DashboardHandler) GetMeetingsByStatus(c *gin.Context) {
	filter := h.parseFilter(c)
	
	statusCounts, err := h.dashboardService.GetMeetingsByStatus(filter)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve meetings by status")
		return
	}
	
	utils.SuccessResponse(c, "Meetings by status retrieved successfully", statusCounts)
}

func (h *DashboardHandler) GetMeetingsByMonth(c *gin.Context) {
	filter := h.parseFilter(c)
	
	monthlyCounts, err := h.dashboardService.GetMeetingsByMonth(filter)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve meetings by month")
		return
	}
	
	utils.SuccessResponse(c, "Meetings by month retrieved successfully", monthlyCounts)
}

func (h *DashboardHandler) GetTopActiveUsers(c *gin.Context) {
	filter := h.parseFilter(c)
	
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	
	users, err := h.dashboardService.GetTopActiveUsers(filter, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve top active users")
		return
	}
	
	utils.SuccessResponse(c, "Top active users retrieved successfully", users)
}

func (h *DashboardHandler) GetRecentMeetings(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	
	meetings, err := h.dashboardService.GetRecentMeetings(limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve recent meetings")
		return
	}
	
	utils.SuccessResponse(c, "Recent meetings retrieved successfully", meetings)
}

func (h *DashboardHandler) GetUpcomingMeetings7d(c *gin.Context) {
	meetings, err := h.dashboardService.GetUpcomingMeetings7d()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve upcoming meetings")
		return
	}
	
	utils.SuccessResponse(c, "Upcoming meetings (7 days) retrieved successfully", meetings)
}

func (h *DashboardHandler) parseFilter(c *gin.Context) models.DashboardFilter {
	var filter models.DashboardFilter
	
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filter.StartDate = &startDate
		}
	}
	
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filter.EndDate = &endDate
		}
	}
	
	if roomIDStr := c.Query("room_id"); roomIDStr != "" {
		if roomID, err := strconv.ParseUint(roomIDStr, 10, 32); err == nil {
			roomIDUint := uint(roomID)
			filter.RoomID = &roomIDUint
		}
	}
	
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			userIDUint := uint(userID)
			filter.UserID = &userIDUint
		}
	}
	
	return filter
}