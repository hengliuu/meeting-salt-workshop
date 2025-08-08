package repositories

import (
	"api/internal/models"
	"time"

	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(room *models.Room) (*models.Room, error)
	GetByID(id uint) (*models.Room, error)
	GetAll(offset, limit int) ([]*models.Room, int64, error)
	Update(room *models.Room) (*models.Room, error)
	Delete(id uint) error
	GetActiveRooms() ([]*models.Room, error)
	SearchRooms(query string, offset, limit int) ([]*models.Room, int64, error)
	GetAvailableRooms(startTime, endTime time.Time, capacity *int) ([]*models.Room, error)
	IsRoomAvailable(roomID uint, startTime, endTime time.Time, excludeMeetingID *uint) (bool, error)
}

type RoomFeatureRepository interface {
	Create(feature *models.RoomFeature) (*models.RoomFeature, error)
	GetByID(id uint) (*models.RoomFeature, error)
	GetAll() ([]*models.RoomFeature, error)
	Update(feature *models.RoomFeature) (*models.RoomFeature, error)
	Delete(id uint) error
}

type roomRepository struct {
	db *gorm.DB
}

type roomFeatureRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

func NewRoomFeatureRepository(db *gorm.DB) RoomFeatureRepository {
	return &roomFeatureRepository{db: db}
}

func (r *roomRepository) Create(room *models.Room) (*models.Room, error) {
	if err := r.db.Create(room).Error; err != nil {
		return nil, err
	}
	return r.GetByID(room.ID)
}

func (r *roomRepository) GetByID(id uint) (*models.Room, error) {
	var room models.Room
	if err := r.db.Preload("Features").Preload("Meetings").First(&room, id).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) GetAll(offset, limit int) ([]*models.Room, int64, error) {
	var rooms []*models.Room
	var total int64

	if err := r.db.Model(&models.Room{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("Features").Offset(offset).Limit(limit).Find(&rooms).Error; err != nil {
		return nil, 0, err
	}

	return rooms, total, nil
}

func (r *roomRepository) Update(room *models.Room) (*models.Room, error) {
	if err := r.db.Save(room).Error; err != nil {
		return nil, err
	}
	return r.GetByID(room.ID)
}

func (r *roomRepository) Delete(id uint) error {
	return r.db.Delete(&models.Room{}, id).Error
}

func (r *roomRepository) GetActiveRooms() ([]*models.Room, error) {
	var rooms []*models.Room
	if err := r.db.Preload("Features").Where("is_active = ?", true).Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *roomRepository) SearchRooms(query string, offset, limit int) ([]*models.Room, int64, error) {
	var rooms []*models.Room
	var total int64

	searchQuery := "%" + query + "%"
	dbQuery := r.db.Model(&models.Room{}).Where(
		"name LIKE ? OR description LIKE ? OR location LIKE ?",
		searchQuery, searchQuery, searchQuery,
	)

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := dbQuery.Preload("Features").Offset(offset).Limit(limit).Find(&rooms).Error; err != nil {
		return nil, 0, err
	}

	return rooms, total, nil
}

func (r *roomRepository) GetAvailableRooms(startTime, endTime time.Time, capacity *int) ([]*models.Room, error) {
	query := r.db.Where("is_active = ?", true)
	
	if capacity != nil {
		query = query.Where("capacity >= ?", *capacity)
	}

	query = query.Where("id NOT IN (?)", r.db.Model(&models.Meeting{}).
		Select("room_id").
		Where("start_time < ? AND end_time > ? AND status IN ?", 
			endTime, startTime, []models.MeetingStatus{models.StatusScheduled, models.StatusInProgress}))

	var rooms []*models.Room
	if err := query.Preload("Features").Find(&rooms).Error; err != nil {
		return nil, err
	}

	return rooms, nil
}

func (r *roomRepository) IsRoomAvailable(roomID uint, startTime, endTime time.Time, excludeMeetingID *uint) (bool, error) {
	query := r.db.Model(&models.Meeting{}).
		Where("room_id = ? AND start_time < ? AND end_time > ? AND status IN ?", 
			roomID, endTime, startTime, []models.MeetingStatus{models.StatusScheduled, models.StatusInProgress})

	if excludeMeetingID != nil {
		query = query.Where("id != ?", *excludeMeetingID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil
}

func (r *roomFeatureRepository) Create(feature *models.RoomFeature) (*models.RoomFeature, error) {
	if err := r.db.Create(feature).Error; err != nil {
		return nil, err
	}
	return feature, nil
}

func (r *roomFeatureRepository) GetByID(id uint) (*models.RoomFeature, error) {
	var feature models.RoomFeature
	if err := r.db.First(&feature, id).Error; err != nil {
		return nil, err
	}
	return &feature, nil
}

func (r *roomFeatureRepository) GetAll() ([]*models.RoomFeature, error) {
	var features []*models.RoomFeature
	if err := r.db.Find(&features).Error; err != nil {
		return nil, err
	}
	return features, nil
}

func (r *roomFeatureRepository) Update(feature *models.RoomFeature) (*models.RoomFeature, error) {
	if err := r.db.Save(feature).Error; err != nil {
		return nil, err
	}
	return feature, nil
}

func (r *roomFeatureRepository) Delete(id uint) error {
	return r.db.Delete(&models.RoomFeature{}, id).Error
}