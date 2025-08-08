# Meeting Salt Workshop API

A comprehensive meeting management system built with Go, Gin, GORM, and MySQL with Microsoft OAuth authentication.

## Features

- **Authentication**: Microsoft OAuth 2.0 integration with JWT tokens
- **User Management**: Role-based access control (Admin, Manager, Employee)
- **Room Management**: Meeting room booking with features and availability checking
- **Meeting Management**: Full CRUD operations with attendee management
- **Dashboard**: Analytics and reporting with usage statistics
- **RESTful API**: Clean API design with proper HTTP status codes
- **Database**: MySQL with GORM ORM for data persistence

## Quick Start

### Prerequisites

- Go 1.21 or higher
- MySQL 8.0 or higher
- Microsoft Azure app registration for OAuth

### Installation

1. Clone the repository
2. Copy the environment file:
   ```bash
   cp .env.example .env
   ```

3. Update `.env` with your configurations:
   - Database credentials
   - Microsoft OAuth credentials
   - JWT secret

4. Install dependencies:
   ```bash
   go mod tidy
   ```

5. Create MySQL database:
   ```sql
   CREATE DATABASE meeting_salt_db;
   ```

6. Run the application:
   ```bash
   go run cmd/server/main.go
   ```

The server will start on `http://localhost:8080`

## API Documentation

### Authentication Endpoints

- `GET /api/v1/auth/login` - Get Microsoft OAuth login URL
- `GET /api/v1/auth/callback` - OAuth callback handler
- `POST /api/v1/auth/refresh` - Refresh JWT token
- `POST /api/v1/auth/logout` - Logout user
- `GET /api/v1/auth/me` - Get current user info

### User Endpoints

- `GET /api/v1/users` - Get all users (paginated)
- `POST /api/v1/users` - Create user (Admin only)
- `GET /api/v1/users/active` - Get active users
- `GET /api/v1/users/search?q=query` - Search users
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Deactivate user (Admin only)

### Room Endpoints

- `GET /api/v1/rooms` - Get all rooms (paginated)
- `POST /api/v1/rooms` - Create room (Manager+)
- `GET /api/v1/rooms/active` - Get active rooms
- `GET /api/v1/rooms/available` - Get available rooms
- `GET /api/v1/rooms/search?q=query` - Search rooms
- `GET /api/v1/rooms/:id` - Get room by ID
- `PUT /api/v1/rooms/:id` - Update room (Manager+)
- `DELETE /api/v1/rooms/:id` - Deactivate room (Manager+)

### Meeting Endpoints

- `GET /api/v1/meetings` - Get all meetings (paginated)
- `POST /api/v1/meetings` - Create meeting
- `GET /api/v1/meetings/filter` - Filter meetings
- `GET /api/v1/meetings/upcoming` - Get upcoming meetings
- `GET /api/v1/meetings/:id` - Get meeting by ID
- `PUT /api/v1/meetings/:id` - Update meeting
- `DELETE /api/v1/meetings/:id` - Delete meeting
- `POST /api/v1/meetings/:id/start` - Start meeting
- `POST /api/v1/meetings/:id/complete` - Complete meeting
- `POST /api/v1/meetings/:id/cancel` - Cancel meeting

### Dashboard Endpoints

- `GET /api/v1/dashboard/stats` - Get dashboard statistics
- `GET /api/v1/dashboard/room-utilization` - Get room utilization data
- `GET /api/v1/dashboard/meetings-by-status` - Get meetings by status
- `GET /api/v1/dashboard/meetings-by-month` - Get monthly meeting counts
- `GET /api/v1/dashboard/top-active-users` - Get most active users

## Project Structure

```
api/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   ├── config.go           # Configuration management
│   │   └── database.go         # Database connection
│   ├── models/
│   │   ├── user.go            # User models
│   │   ├── room.go            # Room models
│   │   ├── meeting.go         # Meeting models
│   │   └── dashboard.go       # Dashboard models
│   ├── handlers/
│   │   ├── auth.go            # Authentication handlers
│   │   ├── users.go           # User handlers
│   │   ├── rooms.go           # Room handlers
│   │   ├── meetings.go        # Meeting handlers
│   │   └── dashboard.go       # Dashboard handlers
│   ├── services/
│   │   ├── auth.go            # Authentication service
│   │   ├── user.go            # User service
│   │   ├── room.go            # Room service
│   │   ├── meeting.go         # Meeting service
│   │   └── dashboard.go       # Dashboard service
│   ├── repositories/
│   │   ├── user.go            # User repository
│   │   ├── room.go            # Room repository
│   │   ├── meeting.go         # Meeting repository
│   │   └── dashboard.go       # Dashboard repository
│   ├── middleware/
│   │   ├── auth.go            # Authentication middleware
│   │   ├── cors.go            # CORS middleware
│   │   └── logger.go          # Logging middleware
│   └── utils/
│       ├── jwt.go             # JWT utilities
│       ├── response.go        # API response utilities
│       ├── validation.go      # Validation utilities
│       └── pagination.go      # Pagination utilities
├── .env.example               # Environment variables template
├── go.mod                     # Go modules
├── go.sum                     # Go modules checksum
├── CLAUDE.md                  # Documentation
└── README.md                  # This file
```

## Authentication Flow

1. Frontend redirects to `/api/v1/auth/login`
2. API returns Microsoft OAuth URL
3. User authenticates with Microsoft
4. Microsoft redirects to `/api/v1/auth/callback`
5. API exchanges code for user info and creates/updates user
6. API returns JWT token
7. Frontend stores token and includes in `Authorization: Bearer <token>` header

## Database Schema

The system uses the following main entities:

- **Users**: Store user information from Microsoft OAuth
- **Rooms**: Meeting rooms with capacity and features
- **RoomFeatures**: Features that rooms can have (projector, whiteboard, etc.)
- **Meetings**: Meeting bookings with organizer and attendees
- **MeetingAttendees**: Many-to-many relationship between meetings and users

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| DB_HOST | MySQL host | localhost |
| DB_PORT | MySQL port | 3306 |
| DB_USER | MySQL username | root |
| DB_PASSWORD | MySQL password | - |
| DB_NAME | MySQL database name | meeting_salt_db |
| SERVER_HOST | Server host | localhost |
| SERVER_PORT | Server port | 8080 |
| MICROSOFT_CLIENT_ID | Microsoft OAuth client ID | - |
| MICROSOFT_CLIENT_SECRET | Microsoft OAuth client secret | - |
| MICROSOFT_REDIRECT_URL | OAuth redirect URL | http://localhost:8080/api/v1/auth/callback |
| JWT_SECRET | JWT signing secret | - |
| GIN_MODE | Gin mode (debug/release) | debug |

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o bin/server cmd/server/main.go
```

### Docker Support (Optional)
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
CMD ["./server"]
```

## Microsoft OAuth Setup

1. Go to Azure Portal > App registrations
2. Create new registration
3. Add redirect URI: `http://localhost:8080/api/v1/auth/callback`
4. Copy Application (client) ID and create client secret
5. Add required permissions: `openid`, `profile`, `email`

## API Response Format

All API endpoints return responses in this format:

```json
{
  "success": true,
  "message": "Operation successful",
  "data": {},
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

## Error Handling

Errors are returned with appropriate HTTP status codes and descriptive messages:

```json
{
  "success": false,
  "message": "Validation failed",
  "error": "Invalid input data",
  "errors": [
    {
      "field": "email",
      "tag": "required",
      "message": "email is required"
    }
  ]
}
```

## Security Features

- JWT authentication with expiration
- Role-based authorization
- Input validation and sanitization
- CORS protection
- Request ID tracking
- Structured logging

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License.