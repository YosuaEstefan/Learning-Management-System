package services

import (
	"LMS/models"
	"LMS/repositories"
	"errors"
	"path/filepath"
	"time"
)

// MaterialService handles material business logic
type MaterialService struct {
	MaterialRepo *repositories.MaterialRepository
	CourseRepo   *repositories.CourseRepository
}

// NewMaterialService creates a new material service
func NewMaterialService(materialRepo *repositories.MaterialRepository, courseRepo *repositories.CourseRepository) *MaterialService {
	return &MaterialService{
		MaterialRepo: materialRepo,
		CourseRepo:   courseRepo,
	}
}

// CreateMaterial creates a new material
func (s *MaterialService) CreateMaterial(material *models.Material, filePath string) error {
	// Verify the course exists
	_, err := s.CourseRepo.FindByID(material.CourseID)
	if err != nil {
		return errors.New("course not found")
	}

	// Set the file path and uploaded date
	material.FilePath = filePath
	material.UploadedAt = time.Now()

	// Create the material
	return s.MaterialRepo.Create(material)
}

// GetMaterialByID gets a material by ID
func (s *MaterialService) GetMaterialByID(id uint) (*models.Material, error) {
	return s.MaterialRepo.FindByID(id)
}

// GetMaterialsByCourse gets materials by course ID
func (s *MaterialService) GetMaterialsByCourse(courseID uint) ([]models.Material, error) {
	return s.MaterialRepo.FindByCourse(courseID)
}

// UpdateMaterial updates a material
func (s *MaterialService) UpdateMaterial(material *models.Material, filePath string) error {
	// Verify the material exists
	existingMaterial, err := s.MaterialRepo.FindByID(material.ID)
	if err != nil {
		return err
	}

	// Update only allowed fields
	existingMaterial.Title = material.Title

	// If a new file is provided, update the file path
	if filePath != "" {
		existingMaterial.FilePath = filePath
	}

	return s.MaterialRepo.Update(existingMaterial)
}

// DeleteMaterial deletes a material
func (s *MaterialService) DeleteMaterial(id uint) error {
	return s.MaterialRepo.Delete(id)
}

// GetFileExtension gets the file extension from a file path
func (s *MaterialService) GetFileExtension(filename string) string {
	return filepath.Ext(filename)
}
