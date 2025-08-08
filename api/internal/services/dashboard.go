package services

import (
	"api/internal/models"
	"api/internal/repositories"
)

type DashboardService struct {
	dashboardRepo repositories.DashboardRepository
}

func NewDashboardService(dashboardRepo repositories.DashboardRepository) *DashboardService {
	return &DashboardService{
		dashboardRepo: dashboardRepo,
	}
}

func (s *DashboardService) GetDashboardStats(filter models.DashboardFilter) (*models.DashboardStats, error) {
	return s.dashboardRepo.GetDashboardStats(filter)
}

func (s *DashboardService) GetRoomUtilization(filter models.DashboardFilter) ([]models.RoomUtilizationData, error) {
	return s.dashboardRepo.GetRoomUtilization(filter)
}

func (s *DashboardService) GetMeetingsByStatus(filter models.DashboardFilter) ([]models.MeetingStatusCount, error) {
	return s.dashboardRepo.GetMeetingsByStatus(filter)
}

func (s *DashboardService) GetMeetingsByMonth(filter models.DashboardFilter) ([]models.MeetingMonthlyCount, error) {
	return s.dashboardRepo.GetMeetingsByMonth(filter)
}

func (s *DashboardService) GetTopActiveUsers(filter models.DashboardFilter, limit int) ([]models.UserActivityData, error) {
	return s.dashboardRepo.GetTopActiveUsers(filter, limit)
}

func (s *DashboardService) GetRecentMeetings(limit int) ([]models.Meeting, error) {
	return s.dashboardRepo.GetRecentMeetings(limit)
}

func (s *DashboardService) GetUpcomingMeetings7d() ([]models.Meeting, error) {
	return s.dashboardRepo.GetUpcomingMeetings7d()
}