package repositories

import (
	"api/internal/models"
	"time"

	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetDashboardStats(filter models.DashboardFilter) (*models.DashboardStats, error)
	GetRoomUtilization(filter models.DashboardFilter) ([]models.RoomUtilizationData, error)
	GetMeetingsByStatus(filter models.DashboardFilter) ([]models.MeetingStatusCount, error)
	GetMeetingsByMonth(filter models.DashboardFilter) ([]models.MeetingMonthlyCount, error)
	GetTopActiveUsers(filter models.DashboardFilter, limit int) ([]models.UserActivityData, error)
	GetRecentMeetings(limit int) ([]models.Meeting, error)
	GetUpcomingMeetings7d() ([]models.Meeting, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetDashboardStats(filter models.DashboardFilter) (*models.DashboardStats, error) {
	stats := &models.DashboardStats{}

	if err := r.db.Model(&models.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&models.User{}).Where("is_active = ?", true).Count(&stats.ActiveUsers).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&models.Room{}).Count(&stats.TotalRooms).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&models.Room{}).Where("is_active = ?", true).Count(&stats.ActiveRooms).Error; err != nil {
		return nil, err
	}

	meetingQuery := r.db.Model(&models.Meeting{})
	if filter.StartDate != nil {
		meetingQuery = meetingQuery.Where("start_time >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		meetingQuery = meetingQuery.Where("end_time <= ?", *filter.EndDate)
	}
	if filter.RoomID != nil {
		meetingQuery = meetingQuery.Where("room_id = ?", *filter.RoomID)
	}
	if filter.UserID != nil {
		meetingQuery = meetingQuery.Where("organizer_id = ? OR id IN (SELECT meeting_id FROM meeting_attendees WHERE user_id = ?)", 
			*filter.UserID, *filter.UserID)
	}

	if err := meetingQuery.Count(&stats.TotalMeetings).Error; err != nil {
		return nil, err
	}

	if err := meetingQuery.Where("status = ?", models.StatusScheduled).Count(&stats.UpcomingMeetings).Error; err != nil {
		return nil, err
	}

	if err := meetingQuery.Where("status = ?", models.StatusCompleted).Count(&stats.CompletedMeetings).Error; err != nil {
		return nil, err
	}

	if err := meetingQuery.Where("status = ?", models.StatusCancelled).Count(&stats.CancelledMeetings).Error; err != nil {
		return nil, err
	}

	roomUtilization, err := r.GetRoomUtilization(filter)
	if err != nil {
		return nil, err
	}
	stats.RoomUtilization = roomUtilization

	meetingsByStatus, err := r.GetMeetingsByStatus(filter)
	if err != nil {
		return nil, err
	}
	stats.MeetingsByStatus = meetingsByStatus

	meetingsByMonth, err := r.GetMeetingsByMonth(filter)
	if err != nil {
		return nil, err
	}
	stats.MeetingsByMonth = meetingsByMonth

	topActiveUsers, err := r.GetTopActiveUsers(filter, 10)
	if err != nil {
		return nil, err
	}
	stats.TopActiveUsers = topActiveUsers

	recentMeetings, err := r.GetRecentMeetings(10)
	if err != nil {
		return nil, err
	}
	stats.RecentMeetings = recentMeetings

	upcomingMeetings7d, err := r.GetUpcomingMeetings7d()
	if err != nil {
		return nil, err
	}
	stats.UpcomingMeetings7d = upcomingMeetings7d

	return stats, nil
}

func (r *dashboardRepository) GetRoomUtilization(filter models.DashboardFilter) ([]models.RoomUtilizationData, error) {
	query := `
		SELECT 
			r.id as room_id,
			r.name as room_name,
			COUNT(m.id) as total_bookings,
			COALESCE(SUM(TIMESTAMPDIFF(MINUTE, m.start_time, m.end_time)) / 60.0, 0) as total_hours,
			COALESCE(SUM(TIMESTAMPDIFF(MINUTE, m.start_time, m.end_time)) / 60.0 / (24 * 30) * 100, 0) as utilization_rate
		FROM rooms r
		LEFT JOIN meetings m ON r.id = m.room_id AND m.status IN ('scheduled', 'completed')
	`

	var conditions []string
	var args []interface{}

	if filter.StartDate != nil {
		conditions = append(conditions, "m.start_time >= ?")
		args = append(args, *filter.StartDate)
	}
	if filter.EndDate != nil {
		conditions = append(conditions, "m.end_time <= ?")
		args = append(args, *filter.EndDate)
	}
	if filter.RoomID != nil {
		conditions = append(conditions, "r.id = ?")
		args = append(args, *filter.RoomID)
	}

	if len(conditions) > 0 {
		query += " WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			query += " AND " + conditions[i]
		}
	}

	query += " GROUP BY r.id, r.name ORDER BY total_bookings DESC"

	var utilization []models.RoomUtilizationData
	if err := r.db.Raw(query, args...).Scan(&utilization).Error; err != nil {
		return nil, err
	}

	return utilization, nil
}

func (r *dashboardRepository) GetMeetingsByStatus(filter models.DashboardFilter) ([]models.MeetingStatusCount, error) {
	query := r.db.Model(&models.Meeting{}).Select("status, COUNT(*) as count")

	if filter.StartDate != nil {
		query = query.Where("start_time >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("end_time <= ?", *filter.EndDate)
	}
	if filter.RoomID != nil {
		query = query.Where("room_id = ?", *filter.RoomID)
	}
	if filter.UserID != nil {
		query = query.Where("organizer_id = ? OR id IN (SELECT meeting_id FROM meeting_attendees WHERE user_id = ?)", 
			*filter.UserID, *filter.UserID)
	}

	var statusCounts []models.MeetingStatusCount
	if err := query.Group("status").Scan(&statusCounts).Error; err != nil {
		return nil, err
	}

	return statusCounts, nil
}

func (r *dashboardRepository) GetMeetingsByMonth(filter models.DashboardFilter) ([]models.MeetingMonthlyCount, error) {
	query := r.db.Model(&models.Meeting{}).
		Select("DATE_FORMAT(start_time, '%Y-%m') as month, COUNT(*) as count")

	if filter.StartDate != nil {
		query = query.Where("start_time >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("end_time <= ?", *filter.EndDate)
	}
	if filter.RoomID != nil {
		query = query.Where("room_id = ?", *filter.RoomID)
	}
	if filter.UserID != nil {
		query = query.Where("organizer_id = ? OR id IN (SELECT meeting_id FROM meeting_attendees WHERE user_id = ?)", 
			*filter.UserID, *filter.UserID)
	}

	var monthlyCounts []models.MeetingMonthlyCount
	if err := query.Group("month").Order("month ASC").Scan(&monthlyCounts).Error; err != nil {
		return nil, err
	}

	return monthlyCounts, nil
}

func (r *dashboardRepository) GetTopActiveUsers(filter models.DashboardFilter, limit int) ([]models.UserActivityData, error) {
	query := `
		SELECT 
			u.id as user_id,
			CONCAT(u.first_name, ' ', u.last_name) as user_name,
			u.email as user_email,
			COUNT(DISTINCT m1.id) as organized_meetings,
			COUNT(DISTINCT ma.meeting_id) as attended_meetings,
			(COUNT(DISTINCT m1.id) + COUNT(DISTINCT ma.meeting_id)) as total_meetings
		FROM users u
		LEFT JOIN meetings m1 ON u.id = m1.organizer_id
		LEFT JOIN meeting_attendees ma ON u.id = ma.user_id
		LEFT JOIN meetings m2 ON ma.meeting_id = m2.id
	`

	var conditions []string
	var args []interface{}

	if filter.StartDate != nil {
		conditions = append(conditions, "(m1.start_time >= ? OR m2.start_time >= ?)")
		args = append(args, *filter.StartDate, *filter.StartDate)
	}
	if filter.EndDate != nil {
		conditions = append(conditions, "(m1.end_time <= ? OR m2.end_time <= ?)")
		args = append(args, *filter.EndDate, *filter.EndDate)
	}
	if filter.RoomID != nil {
		conditions = append(conditions, "(m1.room_id = ? OR m2.room_id = ?)")
		args = append(args, *filter.RoomID, *filter.RoomID)
	}

	if len(conditions) > 0 {
		query += " WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			query += " AND " + conditions[i]
		}
	}

	query += " GROUP BY u.id, u.first_name, u.last_name, u.email ORDER BY total_meetings DESC LIMIT ?"
	args = append(args, limit)

	var userActivity []models.UserActivityData
	if err := r.db.Raw(query, args...).Scan(&userActivity).Error; err != nil {
		return nil, err
	}

	return userActivity, nil
}

func (r *dashboardRepository) GetRecentMeetings(limit int) ([]models.Meeting, error) {
	var meetings []models.Meeting
	if err := r.db.Preload("Organizer").Preload("Room").
		Order("created_at DESC").Limit(limit).Find(&meetings).Error; err != nil {
		return nil, err
	}

	return meetings, nil
}

func (r *dashboardRepository) GetUpcomingMeetings7d() ([]models.Meeting, error) {
	startTime := time.Now()
	endTime := startTime.AddDate(0, 0, 7)

	var meetings []models.Meeting
	if err := r.db.Preload("Organizer").Preload("Room").
		Where("start_time >= ? AND start_time <= ? AND status = ?", 
			startTime, endTime, models.StatusScheduled).
		Order("start_time ASC").Find(&meetings).Error; err != nil {
		return nil, err
	}

	return meetings, nil
}