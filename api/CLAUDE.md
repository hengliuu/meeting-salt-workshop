# Meeting Salt Workshop API

## Project Overview
Backend system for a meeting management platform built with Go, Gin framework, GORM ORM, and MySQL database. The system includes comprehensive authentication via Microsoft OAuth and manages meetings, rooms, users, and dashboard analytics.

## Tech Stack
- **Framework**: Gin (HTTP web framework)
- **ORM**: GORM (Go ORM library)
- **Database**: MySQL 8.0+
- **Authentication**: Microsoft OAuth 2.0
- **Language**: Go 1.21+

## System Domains

### 1. Authentication
- Microsoft OAuth 2.0 integration
- JWT token management
- Role-based access control (RBAC)
- Session management

### 2. Users
- User profile management
- Role assignment (Admin, Manager, Employee)
- User preferences and settings

### 3. Rooms
- Meeting room inventory
- Room capacity and features
- Availability management
- Booking restrictions

### 4. Meetings
- Meeting scheduling and management
- Attendee management
- Meeting status tracking
- Recurring meeting support

### 5. Dashboard
- Analytics and reporting
- Usage statistics
- Real-time metrics
- Data visualization endpoints

## Database Configuration

### MySQL Setup
```go
// config/database.go
package config

import (
    "fmt"
    "log"
    "os"
    
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

type Database struct {
    DB *gorm.DB
}

func NewDatabase() (*Database, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"),
    )
    
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %v", err)
    }
    
    return &Database{DB: db}, nil
}
```

### Environment Variables
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=meeting_salt_db

MICROSOFT_CLIENT_ID=your_client_id
MICROSOFT_CLIENT_SECRET=your_client_secret
MICROSOFT_REDIRECT_URL=http://localhost:8080/auth/callback

JWT_SECRET=your_jwt_secret
```

## Domain Models

### User Model
```go
// models/user.go
package models

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID              uint           `json:"id" gorm:"primaryKey"`
    MicrosoftID     string         `json:"microsoft_id" gorm:"uniqueIndex;not null"`
    Email           string         `json:"email" gorm:"uniqueIndex;not null"`
    FirstName       string         `json:"first_name" gorm:"not null"`
    LastName        string         `json:"last_name" gorm:"not null"`
    DisplayName     string         `json:"display_name"`
    ProfilePicture  string         `json:"profile_picture"`
    Role            UserRole       `json:"role" gorm:"default:'employee'"`
    IsActive        bool           `json:"is_active" gorm:"default:true"`
    LastLogin       *time.Time     `json:"last_login"`
    CreatedAt       time.Time      `json:"created_at"`
    UpdatedAt       time.Time      `json:"updated_at"`
    DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
    
    // Relationships
    OrganizedMeetings []Meeting `json:"organized_meetings" gorm:"foreignKey:OrganizerID"`
    AttendedMeetings  []Meeting `json:"attended_meetings" gorm:"many2many:meeting_attendees;"`
}

type UserRole string

const (
    RoleAdmin    UserRole = "admin"
    RoleManager  UserRole = "manager"
    RoleEmployee UserRole = "employee"
)
```

### Room Model
```go
// models/room.go
package models

import (
    "time"
    "gorm.io/gorm"
)

type Room struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Name        string         `json:"name" gorm:"not null"`
    Description string         `json:"description"`
    Capacity    int            `json:"capacity" gorm:"not null"`
    Location    string         `json:"location"`
    Features    []RoomFeature  `json:"features" gorm:"many2many:room_room_features;"`
    IsActive    bool           `json:"is_active" gorm:"default:true"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
    
    // Relationships
    Meetings []Meeting `json:"meetings" gorm:"foreignKey:RoomID"`
}

type RoomFeature struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Name        string         `json:"name" gorm:"uniqueIndex;not null"`
    Description string         `json:"description"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
```

### Meeting Model
```go
// models/meeting.go
package models

import (
    "time"
    "gorm.io/gorm"
)

type Meeting struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Title       string         `json:"title" gorm:"not null"`
    Description string         `json:"description"`
    StartTime   time.Time      `json:"start_time" gorm:"not null"`
    EndTime     time.Time      `json:"end_time" gorm:"not null"`
    Status      MeetingStatus  `json:"status" gorm:"default:'scheduled'"`
    IsRecurring bool           `json:"is_recurring" gorm:"default:false"`
    RecurrencePattern string   `json:"recurrence_pattern"`
    
    // Foreign Keys
    OrganizerID uint `json:"organizer_id" gorm:"not null"`
    RoomID      uint `json:"room_id" gorm:"not null"`
    
    // Timestamps
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
    
    // Relationships
    Organizer User   `json:"organizer" gorm:"foreignKey:OrganizerID"`
    Room      Room   `json:"room" gorm:"foreignKey:RoomID"`
    Attendees []User `json:"attendees" gorm:"many2many:meeting_attendees;"`
}

type MeetingStatus string

const (
    StatusScheduled MeetingStatus = "scheduled"
    StatusInProgress MeetingStatus = "in_progress"
    StatusCompleted MeetingStatus = "completed"
    StatusCancelled MeetingStatus = "cancelled"
)
```

## Microsoft OAuth Integration

### OAuth Configuration
```go
// auth/microsoft.go
package auth

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/microsoft"
)

type MicrosoftAuth struct {
    Config *oauth2.Config
}

func NewMicrosoftAuth() *MicrosoftAuth {
    config := &oauth2.Config{
        ClientID:     os.Getenv("MICROSOFT_CLIENT_ID"),
        ClientSecret: os.Getenv("MICROSOFT_CLIENT_SECRET"),
        RedirectURL:  os.Getenv("MICROSOFT_REDIRECT_URL"),
        Scopes:       []string{"openid", "profile", "email"},
        Endpoint:     microsoft.AzureADEndpoint(""),
    }
    
    return &MicrosoftAuth{Config: config}
}

type MicrosoftUser struct {
    ID          string `json:"id"`
    Email       string `json:"mail"`
    FirstName   string `json:"givenName"`
    LastName    string `json:"surname"`
    DisplayName string `json:"displayName"`
}

func (m *MicrosoftAuth) GetUserInfo(token *oauth2.Token) (*MicrosoftUser, error) {
    client := m.Config.Client(context.Background(), token)
    
    resp, err := client.Get("https://graph.microsoft.com/v1.0/me")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var user MicrosoftUser
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        return nil, err
    }
    
    return &user, nil
}
```

## Project Structure
```
api/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   ├── config.go
│   │   └── database.go
│   ├── models/
│   │   ├── user.go
│   │   ├── room.go
│   │   ├── meeting.go
│   │   └── dashboard.go
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── users.go
│   │   ├── rooms.go
│   │   ├── meetings.go
│   │   └── dashboard.go
│   ├── services/
│   │   ├── auth.go
│   │   ├── user.go
│   │   ├── room.go
│   │   ├── meeting.go
│   │   └── dashboard.go
│   ├── repositories/
│   │   ├── user.go
│   │   ├── room.go
│   │   ├── meeting.go
│   │   └── dashboard.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logger.go
│   └── utils/
│       ├── jwt.go
│       ├── validation.go
│       └── response.go
├── migrations/
├── docs/
├── .env
├── .gitignore
├── go.mod
├── go.sum
└── CLAUDE.md
```

## Best Practices

### 1. Error Handling
```go
// utils/response.go
package utils

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type APIResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func SuccessResponse(c *gin.Context, message string, data interface{}) {
    c.JSON(http.StatusOK, APIResponse{
        Success: true,
        Message: message,
        Data:    data,
    })
}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
    c.JSON(statusCode, APIResponse{
        Success: false,
        Message: message,
        Error:   message,
    })
}
```

### 2. Validation
```go
// utils/validation.go
package utils

import (
    "github.com/go-playground/validator/v10"
)

type ValidationError struct {
    Field   string `json:"field"`
    Tag     string `json:"tag"`
    Message string `json:"message"`
}

func ValidateStruct(s interface{}) []ValidationError {
    var errors []ValidationError
    validate := validator.New()
    
    err := validate.Struct(s)
    if err != nil {
        for _, err := range err.(validator.ValidationErrors) {
            errors = append(errors, ValidationError{
                Field:   err.Field(),
                Tag:     err.Tag(),
                Message: getValidationMessage(err),
            })
        }
    }
    
    return errors
}
```

### 3. Database Migrations
- Use GORM AutoMigrate for development
- Use proper SQL migrations for production
- Always backup before migrations
- Version control migration files

### 4. Security Best Practices
- Use HTTPS in production
- Implement rate limiting
- Validate all inputs
- Use parameterized queries
- Implement proper CORS
- Store secrets in environment variables
- Use JWT with appropriate expiration

### 5. API Design
- Follow RESTful conventions
- Use consistent response formats
- Implement proper pagination
- Use appropriate HTTP status codes
- Document API endpoints

## Development Commands
```bash
# Run the application
go run cmd/server/main.go

# Run tests
go test ./...

# Build for production
go build -o bin/server cmd/server/main.go

# Database migrations
go run migrations/migrate.go
```