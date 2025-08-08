package handlers

import (
	"api/internal/models"
	"api/internal/services"
	"api/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	roomService *services.RoomService
}

func NewRoomHandler(roomService *services.RoomService) *RoomHandler {
	return &RoomHandler{
		roomService: roomService,
	}
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var req models.CreateRoomRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request format")
		return
	}
	
	if validationErrors := utils.ValidateStruct(&req); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}
	
	room, err := h.roomService.CreateRoom(&req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.CreatedResponse(c, "Room created successfully", room)
}

func (h *RoomHandler) GetRoom(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid room ID")
		return
	}
	
	room, err := h.roomService.GetRoomByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Room not found")
		return
	}
	
	utils.SuccessResponse(c, "Room retrieved successfully", room)
}

func (h *RoomHandler) GetAllRooms(c *gin.Context) {
	pagination := utils.GetPaginationParams(c)
	
	rooms, meta, err := h.roomService.GetAllRooms(pagination.Page, pagination.Limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve rooms")
		return
	}
	
	utils.PaginatedSuccessResponse(c, "Rooms retrieved successfully", rooms, meta)
}

func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid room ID")
		return
	}
	
	var req models.UpdateRoomRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request format")
		return
	}
	
	room, err := h.roomService.UpdateRoom(uint(id), &req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, "Room updated successfully", room)
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid room ID")
		return
	}
	
	if err := h.roomService.DeleteRoom(uint(id)); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, "Room deactivated successfully", nil)
}

func (h *RoomHandler) GetActiveRooms(c *gin.Context) {
	rooms, err := h.roomService.GetActiveRooms()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve active rooms")
		return
	}
	
	utils.SuccessResponse(c, "Active rooms retrieved successfully", rooms)
}

func (h *RoomHandler) SearchRooms(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.BadRequestResponse(c, "Search query is required")
		return
	}
	
	pagination := utils.GetPaginationParams(c)
	
	rooms, meta, err := h.roomService.SearchRooms(query, pagination.Page, pagination.Limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to search rooms")
		return
	}
	
	utils.PaginatedSuccessResponse(c, "Search results retrieved successfully", rooms, meta)
}

func (h *RoomHandler) GetAvailableRooms(c *gin.Context) {
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")
	capacityStr := c.Query("capacity")
	
	if startTimeStr == "" || endTimeStr == "" {
		utils.BadRequestResponse(c, "start_time and end_time are required")
		return
	}
	
	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid start_time format. Use RFC3339 format")
		return
	}
	
	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid end_time format. Use RFC3339 format")
		return
	}
	
	req := &models.RoomAvailabilityQuery{
		StartTime: startTime,
		EndTime:   endTime,
	}
	
	if capacityStr != "" {
		capacity, err := strconv.Atoi(capacityStr)
		if err != nil {
			utils.BadRequestResponse(c, "Invalid capacity value")
			return
		}
		req.Capacity = &capacity
	}
	
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}
	
	rooms, err := h.roomService.GetAvailableRooms(req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, "Available rooms retrieved successfully", rooms)
}

func (h *RoomHandler) CheckRoomAvailability(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid room ID")
		return
	}
	
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")
	excludeMeetingIDStr := c.Query("exclude_meeting_id")
	
	if startTimeStr == "" || endTimeStr == "" {
		utils.BadRequestResponse(c, "start_time and end_time are required")
		return
	}
	
	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid start_time format. Use RFC3339 format")
		return
	}
	
	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid end_time format. Use RFC3339 format")
		return
	}
	
	var excludeMeetingID *uint
	if excludeMeetingIDStr != "" {
		excludeID, err := strconv.ParseUint(excludeMeetingIDStr, 10, 32)
		if err != nil {
			utils.BadRequestResponse(c, "Invalid exclude_meeting_id value")
			return
		}
		excludeIDUint := uint(excludeID)
		excludeMeetingID = &excludeIDUint
	}
	
	isAvailable, err := h.roomService.IsRoomAvailable(uint(id), startTime, endTime, excludeMeetingID)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, "Room availability checked", gin.H{
		"available": isAvailable,
	})
}

func (h *RoomHandler) CreateFeature(c *gin.Context) {
	var req models.CreateRoomFeatureRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request format")
		return
	}
	
	if validationErrors := utils.ValidateStruct(&req); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}
	
	feature, err := h.roomService.CreateFeature(&req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.CreatedResponse(c, "Room feature created successfully", feature)
}

func (h *RoomHandler) GetFeature(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid feature ID")
		return
	}
	
	feature, err := h.roomService.GetFeatureByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Room feature not found")
		return
	}
	
	utils.SuccessResponse(c, "Room feature retrieved successfully", feature)
}

func (h *RoomHandler) GetAllFeatures(c *gin.Context) {
	features, err := h.roomService.GetAllFeatures()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve room features")
		return
	}
	
	utils.SuccessResponse(c, "Room features retrieved successfully", features)
}

func (h *RoomHandler) UpdateFeature(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid feature ID")
		return
	}
	
	var req models.UpdateRoomFeatureRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request format")
		return
	}
	
	feature, err := h.roomService.UpdateFeature(uint(id), &req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, "Room feature updated successfully", feature)
}

func (h *RoomHandler) DeleteFeature(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid feature ID")
		return
	}
	
	if err := h.roomService.DeleteFeature(uint(id)); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, "Room feature deleted successfully", nil)
}