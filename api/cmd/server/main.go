package main

import (
	"api/internal/config"
	"api/internal/handlers"
	"api/internal/middleware"
	"api/internal/models"
	"api/internal/repositories"
	"api/internal/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db, err := config.NewDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := migrateDatabase(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	userRepo := repositories.NewUserRepository(db.DB)
	roomRepo := repositories.NewRoomRepository(db.DB)
	roomFeatureRepo := repositories.NewRoomFeatureRepository(db.DB)
	meetingRepo := repositories.NewMeetingRepository(db.DB)
	dashboardRepo := repositories.NewDashboardRepository(db.DB)

	authService := services.NewAuthService(userRepo, cfg)
	userService := services.NewUserService(userRepo)
	roomService := services.NewRoomService(roomRepo, roomFeatureRepo)
	meetingService := services.NewMeetingService(meetingRepo, roomRepo, userRepo)
	dashboardService := services.NewDashboardService(dashboardRepo)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	roomHandler := handlers.NewRoomHandler(roomService)
	meetingHandler := handlers.NewMeetingHandler(meetingService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	r := setupRouter(cfg, authService, authHandler, userHandler, roomHandler, meetingHandler, dashboardHandler)

	log.Printf("Server starting on http://%s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := r.Run(cfg.Server.Host + ":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRouter(
	cfg *config.Config,
	authService *services.AuthService,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	roomHandler *handlers.RoomHandler,
	meetingHandler *handlers.MeetingHandler,
	dashboardHandler *handlers.DashboardHandler,
) *gin.Engine {
	r := gin.New()

	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.ErrorMiddleware())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "meeting-salt-api",
		})
	})

	api := r.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.GET("/login", authHandler.Login)
		auth.GET("/callback", authHandler.Callback)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", middleware.AuthMiddleware(authService), authHandler.Logout)
		auth.GET("/me", middleware.AuthMiddleware(authService), authHandler.Me)
	}

	users := api.Group("/users")
	users.Use(middleware.AuthMiddleware(authService))
	{
		users.POST("", middleware.RequireRole(models.RoleAdmin), userHandler.CreateUser)
		users.GET("", userHandler.GetAllUsers)
		users.GET("/active", userHandler.GetActiveUsers)
		users.GET("/search", userHandler.SearchUsers)
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id", middleware.RequireOwnerOrRole(models.RoleManager), userHandler.UpdateUser)
		users.DELETE("/:id", middleware.RequireRole(models.RoleAdmin), userHandler.DeleteUser)
		users.PUT("/:id/role", middleware.RequireRole(models.RoleAdmin), userHandler.UpdateUserRole)
		users.POST("/:id/activate", middleware.RequireRole(models.RoleAdmin), userHandler.ActivateUser)
		users.POST("/:id/deactivate", middleware.RequireRole(models.RoleAdmin), userHandler.DeactivateUser)
	}

	rooms := api.Group("/rooms")
	rooms.Use(middleware.AuthMiddleware(authService))
	{
		rooms.POST("", middleware.RequireRole(models.RoleManager), roomHandler.CreateRoom)
		rooms.GET("", roomHandler.GetAllRooms)
		rooms.GET("/active", roomHandler.GetActiveRooms)
		rooms.GET("/available", roomHandler.GetAvailableRooms)
		rooms.GET("/search", roomHandler.SearchRooms)
		rooms.GET("/:id", roomHandler.GetRoom)
		rooms.PUT("/:id", middleware.RequireRole(models.RoleManager), roomHandler.UpdateRoom)
		rooms.DELETE("/:id", middleware.RequireRole(models.RoleManager), roomHandler.DeleteRoom)
		rooms.GET("/:id/availability", roomHandler.CheckRoomAvailability)
	}

	roomFeatures := api.Group("/room-features")
	roomFeatures.Use(middleware.AuthMiddleware(authService))
	{
		roomFeatures.POST("", middleware.RequireRole(models.RoleManager), roomHandler.CreateFeature)
		roomFeatures.GET("", roomHandler.GetAllFeatures)
		roomFeatures.GET("/:id", roomHandler.GetFeature)
		roomFeatures.PUT("/:id", middleware.RequireRole(models.RoleManager), roomHandler.UpdateFeature)
		roomFeatures.DELETE("/:id", middleware.RequireRole(models.RoleManager), roomHandler.DeleteFeature)
	}

	meetings := api.Group("/meetings")
	meetings.Use(middleware.AuthMiddleware(authService))
	{
		meetings.POST("", meetingHandler.CreateMeeting)
		meetings.GET("", meetingHandler.GetAllMeetings)
		meetings.GET("/filter", meetingHandler.GetMeetingsByFilter)
		meetings.GET("/upcoming", meetingHandler.GetUpcomingMeetings)
		meetings.GET("/date-range", meetingHandler.GetMeetingsByDateRange)
		meetings.GET("/:id", meetingHandler.GetMeeting)
		meetings.PUT("/:id", meetingHandler.UpdateMeeting)
		meetings.DELETE("/:id", meetingHandler.DeleteMeeting)
		meetings.POST("/:id/start", meetingHandler.StartMeeting)
		meetings.POST("/:id/complete", meetingHandler.CompleteMeeting)
		meetings.POST("/:id/cancel", meetingHandler.CancelMeeting)
		meetings.GET("/:id/attendees", meetingHandler.GetMeetingAttendees)
		meetings.POST("/:id/attendees", meetingHandler.AddAttendee)
		meetings.DELETE("/:id/attendees/:user_id", meetingHandler.RemoveAttendee)
	}

	dashboard := api.Group("/dashboard")
	dashboard.Use(middleware.AuthMiddleware(authService))
	{
		dashboard.GET("/stats", dashboardHandler.GetDashboardStats)
		dashboard.GET("/room-utilization", dashboardHandler.GetRoomUtilization)
		dashboard.GET("/meetings-by-status", dashboardHandler.GetMeetingsByStatus)
		dashboard.GET("/meetings-by-month", dashboardHandler.GetMeetingsByMonth)
		dashboard.GET("/top-active-users", dashboardHandler.GetTopActiveUsers)
		dashboard.GET("/recent-meetings", dashboardHandler.GetRecentMeetings)
		dashboard.GET("/upcoming-meetings-7d", dashboardHandler.GetUpcomingMeetings7d)
	}

	return r
}

func migrateDatabase(db *config.Database) error {
	return db.DB.AutoMigrate(
		&models.User{},
		&models.Room{},
		&models.RoomFeature{},
		&models.Meeting{},
	)
}