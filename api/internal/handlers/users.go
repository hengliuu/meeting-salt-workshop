package handlers

import (
	"api/internal/middleware"
	"api/internal/models"
	"api/internal/services"
	"api/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup, jwtManager *utils.JWTManager) {
	users := router.Group("/users")
	users.Use(middleware.AuthMiddleware(jwtManager))
	{
		users.GET("", h.GetUsers)
		users.GET("/search", h.SearchUsers)
		users.GET("/active", h.GetActiveUsers)
		users.GET("/:id", h.GetUserByID)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", middleware.RequireAdmin(), h.DeleteUser)
		users.PUT("/:id/role", middleware.RequireAdmin(), h.UpdateUserRole)
		users.PUT("/:id/activate", middleware.RequireAdmin(), h.ActivateUser)
		users.PUT("/:id/deactivate", middleware.RequireAdmin(), h.DeactivateUser)
	}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	paginationParams := utils.GetPaginationParams(c)
	
	users, meta, err := h.userService.GetAllUsers(paginationParams.Page, paginationParams.Limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve users")
		return
	}

	utils.PaginatedSuccessResponse(c, "Users retrieved successfully", users, meta)
}

func (h *UserHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.BadRequestResponse(c, "Search query is required")
		return
	}

	paginationParams := utils.GetPaginationParams(c)
	
	users, meta, err := h.userService.SearchUsers(query, paginationParams.Page, paginationParams.Limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to search users")
		return
	}

	utils.PaginatedSuccessResponse(c, "Users found", users, meta)
}

func (h *UserHandler) GetActiveUsers(c *gin.Context) {
	users, err := h.userService.GetActiveUsers()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve active users")
		return
	}

	utils.SuccessResponse(c, "Active users retrieved successfully", users)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	utils.SuccessResponse(c, "User retrieved successfully", user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}

	currentUserID := middleware.GetUserIDFromContext(c)
	currentUserRole := middleware.GetUserRoleFromContext(c)

	if uint(id) != currentUserID && currentUserRole != "admin" && currentUserRole != "manager" {
		utils.ForbiddenResponse(c, "You can only update your own profile")
		return
	}

	var request models.UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	if request.Role != nil && currentUserRole != "admin" {
		utils.ForbiddenResponse(c, "Only administrators can update user roles")
		return
	}

	user, err := h.userService.UpdateUser(uint(id), &request)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, "User updated successfully", user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, "User deleted successfully", nil)
}

func (h *UserHandler) UpdateUserRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}

	var request struct {
		Role models.UserRole `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	if err := h.userService.UpdateUserRole(uint(id), request.Role); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, "User role updated successfully", nil)
}

func (h *UserHandler) ActivateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}

	if err := h.userService.ActivateUser(uint(id)); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, "User activated successfully", nil)
}

func (h *UserHandler) DeactivateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}

	if err := h.userService.DeactivateUser(uint(id)); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, "User deactivated successfully", nil)
}