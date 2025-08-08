package services

import (
	"api/internal/models"
	"api/internal/repositories"
	"api/internal/utils"
	"errors"
	"time"
)

type RoomService struct {
	roomRepo    repositories.RoomRepository
	featureRepo repositories.RoomFeatureRepository
}

func NewRoomService(roomRepo repositories.RoomRepository, featureRepo repositories.RoomFeatureRepository) *RoomService {
	return &RoomService{
		roomRepo:    roomRepo,
		featureRepo: featureRepo,
	}
}

func (s *RoomService) CreateRoom(req *models.CreateRoomRequest) (*models.Room, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.New("validation failed")
	}

	room := &models.Room{
		Name:        req.Name,
		Description: req.Description,
		Capacity:    req.Capacity,
		Location:    req.Location,
		IsActive:    true,
	}

	createdRoom, err := s.roomRepo.Create(room)
	if err != nil {
		return nil, err
	}

	if len(req.FeatureIDs) > 0 {
		var features []models.RoomFeature
		for _, featureID := range req.FeatureIDs {
			feature, err := s.featureRepo.GetByID(featureID)
			if err != nil {
				continue
			}
			features = append(features, *feature)
		}
		createdRoom.Features = features
		createdRoom, err = s.roomRepo.Update(createdRoom)
		if err != nil {
			return nil, err
		}
	}

	return createdRoom, nil
}

func (s *RoomService) GetRoomByID(id uint) (*models.Room, error) {
	return s.roomRepo.GetByID(id)
}

func (s *RoomService) GetAllRooms(page, limit int) ([]*models.Room, utils.PaginationMeta, error) {
	offset := utils.GetOffset(page, limit)
	rooms, total, err := s.roomRepo.GetAll(offset, limit)
	if err != nil {
		return nil, utils.PaginationMeta{}, err
	}

	meta := utils.CreatePaginationMeta(page, limit, total)
	return rooms, meta, nil
}

func (s *RoomService) UpdateRoom(id uint, req *models.UpdateRoomRequest) (*models.Room, error) {
	room, err := s.roomRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("room not found")
	}

	if req.Name != nil {
		room.Name = *req.Name
	}
	if req.Description != nil {
		room.Description = *req.Description
	}
	if req.Capacity != nil {
		if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
			return nil, errors.New("validation failed")
		}
		room.Capacity = *req.Capacity
	}
	if req.Location != nil {
		room.Location = *req.Location
	}
	if req.IsActive != nil {
		room.IsActive = *req.IsActive
	}

	if req.FeatureIDs != nil {
		var features []models.RoomFeature
		for _, featureID := range req.FeatureIDs {
			feature, err := s.featureRepo.GetByID(featureID)
			if err != nil {
				continue
			}
			features = append(features, *feature)
		}
		room.Features = features
	}

	return s.roomRepo.Update(room)
}

func (s *RoomService) DeleteRoom(id uint) error {
	room, err := s.roomRepo.GetByID(id)
	if err != nil {
		return errors.New("room not found")
	}

	if !room.IsActive {
		return errors.New("room is already inactive")
	}

	room.IsActive = false
	_, err = s.roomRepo.Update(room)
	return err
}

func (s *RoomService) GetActiveRooms() ([]*models.Room, error) {
	return s.roomRepo.GetActiveRooms()
}

func (s *RoomService) SearchRooms(query string, page, limit int) ([]*models.Room, utils.PaginationMeta, error) {
	offset := utils.GetOffset(page, limit)
	rooms, total, err := s.roomRepo.SearchRooms(query, offset, limit)
	if err != nil {
		return nil, utils.PaginationMeta{}, err
	}

	meta := utils.CreatePaginationMeta(page, limit, total)
	return rooms, meta, nil
}

func (s *RoomService) GetAvailableRooms(req *models.RoomAvailabilityQuery) ([]*models.Room, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.New("validation failed")
	}

	if req.EndTime.Before(req.StartTime) {
		return nil, errors.New("end time must be after start time")
	}

	return s.roomRepo.GetAvailableRooms(req.StartTime, req.EndTime, req.Capacity)
}

func (s *RoomService) IsRoomAvailable(roomID uint, startTime, endTime time.Time, excludeMeetingID *uint) (bool, error) {
	if endTime.Before(startTime) {
		return false, errors.New("end time must be after start time")
	}

	return s.roomRepo.IsRoomAvailable(roomID, startTime, endTime, excludeMeetingID)
}

func (s *RoomService) CreateFeature(req *models.CreateRoomFeatureRequest) (*models.RoomFeature, error) {
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, errors.New("validation failed")
	}

	feature := &models.RoomFeature{
		Name:        req.Name,
		Description: req.Description,
	}

	return s.featureRepo.Create(feature)
}

func (s *RoomService) GetFeatureByID(id uint) (*models.RoomFeature, error) {
	return s.featureRepo.GetByID(id)
}

func (s *RoomService) GetAllFeatures() ([]*models.RoomFeature, error) {
	return s.featureRepo.GetAll()
}

func (s *RoomService) UpdateFeature(id uint, req *models.UpdateRoomFeatureRequest) (*models.RoomFeature, error) {
	feature, err := s.featureRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("feature not found")
	}

	if req.Name != nil {
		feature.Name = *req.Name
	}
	if req.Description != nil {
		feature.Description = *req.Description
	}

	return s.featureRepo.Update(feature)
}

func (s *RoomService) DeleteFeature(id uint) error {
	_, err := s.featureRepo.GetByID(id)
	if err != nil {
		return errors.New("feature not found")
	}

	return s.featureRepo.Delete(id)
}