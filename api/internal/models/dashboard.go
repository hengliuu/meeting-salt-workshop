package models

import "time"

type DashboardStats struct {
	TotalUsers         int64                    `json:"total_users"`
	ActiveUsers        int64                    `json:"active_users"`
	TotalRooms         int64                    `json:"total_rooms"`
	ActiveRooms        int64                    `json:"active_rooms"`
	TotalMeetings      int64                    `json:"total_meetings"`
	UpcomingMeetings   int64                    `json:"upcoming_meetings"`
	CompletedMeetings  int64                    `json:"completed_meetings"`
	CancelledMeetings  int64                    `json:"cancelled_meetings"`
	RoomUtilization    []RoomUtilizationData    `json:"room_utilization"`
	MeetingsByStatus   []MeetingStatusCount     `json:"meetings_by_status"`
	MeetingsByMonth    []MeetingMonthlyCount    `json:"meetings_by_month"`
	TopActiveUsers     []UserActivityData       `json:"top_active_users"`
	RecentMeetings     []Meeting                `json:"recent_meetings"`
	UpcomingMeetings7d []Meeting                `json:"upcoming_meetings_7d"`
}

type RoomUtilizationData struct {
	RoomID           uint    `json:"room_id"`
	RoomName         string  `json:"room_name"`
	TotalBookings    int64   `json:"total_bookings"`
	TotalHours       float64 `json:"total_hours"`
	UtilizationRate  float64 `json:"utilization_rate"`
}

type MeetingStatusCount struct {
	Status MeetingStatus `json:"status"`
	Count  int64         `json:"count"`
}

type MeetingMonthlyCount struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}

type UserActivityData struct {
	UserID            uint   `json:"user_id"`
	UserName          string `json:"user_name"`
	UserEmail         string `json:"user_email"`
	OrganizedMeetings int64  `json:"organized_meetings"`
	AttendedMeetings  int64  `json:"attended_meetings"`
	TotalMeetings     int64  `json:"total_meetings"`
}

type DashboardFilter struct {
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	RoomID    *uint      `json:"room_id"`
	UserID    *uint      `json:"user_id"`
}