package services

import (
	"LMS/models"
	"LMS/repositories"
	"errors"
	"path/filepath"
	"time"

)

// MaterialService menangani logika bisnis material
type MaterialService struct {
	MaterialRepo *repositories.MaterialRepository
	CourseRepo   *repositories.CourseRepository
}

// NewMaterialService membuat layanan material baru
func NewMaterialService(materialRepo *repositories.MaterialRepository, courseRepo *repositories.CourseRepository) *MaterialService {
	return &MaterialService{
		MaterialRepo: materialRepo,
		CourseRepo:   courseRepo,
	}
}

// CreateMaterial membuat materi baru
func (s *MaterialService) CreateMaterial(material *models.Material, filePath string) error {
	// Verifikasi keberadaan kursus
	_, err := s.CourseRepo.FindByID(material.CourseID)
	if err != nil {
		return errors.New("course not found")
	}

	// Mengatur jalur file dan tanggal pengunggahan
	material.FilePath = filePath
	material.UploadedAt = time.Now()

	// Buat materi
	return s.MaterialRepo.Create(material)
}

// GetMaterialByID mendapatkan material dengan ID
func (s *MaterialService) GetMaterialByID(id uint) (*models.Material, error) {
	return s.MaterialRepo.FindByID(id)
}

// GetMaterialsByCourse mendapatkan materi berdasarkan ID kursus
func (s *MaterialService) GetMaterialsByCourse(courseID uint) ([]models.Material, error) {
	return s.MaterialRepo.FindByCourse(courseID)
}

// UpdateMaterial memperbarui materi
func (s *MaterialService) UpdateMaterial(material *models.Material, filePath string) error {
	// Verifikasi materi yang ada
	existingMaterial, err := s.MaterialRepo.FindByID(material.ID)
	if err != nil {
		return err
	}

	// Perbarui hanya bidang yang diizinkan
	existingMaterial.Title = material.Title

	// Jika file baru disediakan, perbarui jalur file
	if filePath != "" {
		existingMaterial.FilePath = filePath
	}

	return s.MaterialRepo.Update(existingMaterial)
}

// DeleteMaterial menghapus material
func (s *MaterialService) DeleteMaterial(id uint) error {
	return s.MaterialRepo.Delete(id)
}

// GetFileExtension mendapatkan ekstensi file dari jalur file
func (s *MaterialService) GetFileExtension(filename string) string {
	return filepath.Ext(filename)
}
