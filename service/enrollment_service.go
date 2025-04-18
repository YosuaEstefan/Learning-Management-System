package services

import (
	"LMS/models"
	"LMS/repositories"
)

var EnrollmentServiceInstance EnrollmentService

type EnrollmentService interface {
	Enroll(enrollment *models.Enrollment) error
	GetEnrollment(userID, courseID uint) (*models.Enrollment, error)
	GetEnrollmentsByUser(userID uint) ([]models.Enrollment, error)
}

type enrollmentService struct {
	repo repositories.EnrollmentRepository
}

func NewEnrollmentService(repo repositories.EnrollmentRepository) EnrollmentService {
	return &enrollmentService{repo}
}

func (s *enrollmentService) Enroll(enrollment *models.Enrollment) error {
	return s.repo.Create(enrollment)
}

func (s *enrollmentService) GetEnrollment(userID, courseID uint) (*models.Enrollment, error) {
	return s.repo.GetByUserAndCourse(userID, courseID)
}

func (s *enrollmentService) GetEnrollmentsByUser(userID uint) ([]models.Enrollment, error) {
	return s.repo.GetEnrollmentsByUser(userID)
}
