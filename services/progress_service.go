// services/learning_progress_service.go
package services

import (
	"LMS/models"
	"LMS/repositories"
	"errors"
	"time"
)

// LearningProgressService handles learning progress business logic
type LearningProgressService struct {
	ProgressRepo   *repositories.LearningProgressRepository
	UserRepo       *repositories.UserRepository
	CourseRepo     *repositories.CourseRepository
	EnrollmentRepo *repositories.EnrollmentRepository
	AssignmentRepo *repositories.AssignmentRepository
}

// NewLearningProgressService creates a new learning progress service
func NewLearningProgressService(
	progressRepo *repositories.LearningProgressRepository,
	userRepo *repositories.UserRepository,
	courseRepo *repositories.CourseRepository,
	enrollmentRepo *repositories.EnrollmentRepository,
	assignmentRepo *repositories.AssignmentRepository,
) *LearningProgressService {
	return &LearningProgressService{
		ProgressRepo:   progressRepo,
		UserRepo:       userRepo,
		CourseRepo:     courseRepo,
		EnrollmentRepo: enrollmentRepo,
		AssignmentRepo: assignmentRepo,
	}
}

// CreateOrUpdateGrade allows mentors/admins to grade student's activity
func (s *LearningProgressService) CreateOrUpdateGrade(
	graderId uint,
	studentID uint,
	courseID uint,
	activityType models.ProgressType,
	activityID uint,
	score float64,
	maxScore float64,
	feedback string,
	completed bool,
) error {
	// Verify the grader exists and is mentor or admin
	grader, err := s.UserRepo.FindByID(graderId)
	if err != nil {
		return errors.New("grader not found")
	}

	if grader.Role != models.RoleMentor && grader.Role != models.RoleAdmin {
		return errors.New("only mentors and admins can grade student activities")
	}

	// Verify the student exists
	student, err := s.UserRepo.FindByID(studentID)
	if err != nil {
		return errors.New("student not found")
	}

	if student.Role != models.RoleStudent {
		return errors.New("only students can be graded")
	}

	// Verify the course exists
	course, err := s.CourseRepo.FindByID(courseID)
	if err != nil {
		return errors.New("course not found")
	}

	// For mentors: verify they are assigned to this course
	if grader.Role == models.RoleMentor && course.MentorID != graderId {
		return errors.New("mentors can only grade courses they are assigned to")
	}

	// Verify the student is enrolled in the course
	_, err = s.EnrollmentRepo.FindByUserAndCourse(studentID, courseID)
	if err != nil {
		return errors.New("student is not enrolled in this course")
	}

	// For assignment: verify it exists
	if activityType == models.ProgressTypeAssignment {
		_, err = s.AssignmentRepo.FindByID(activityID)
		if err != nil {
			return errors.New("assignment not found")
		}
	}

	// Check if progress already exists
	now := time.Now()
	existingProgress, err := s.ProgressRepo.FindByActivityAndStudent(
		activityType,
		activityID,
		studentID,
		courseID,
	)

	scoreValue := score
	maxScoreValue := maxScore

	if err == nil && existingProgress != nil {
		// Update existing progress
		existingProgress.Score = &scoreValue
		existingProgress.MaxScore = &maxScoreValue
		existingProgress.Feedback = feedback
		existingProgress.GradedBy = graderId
		existingProgress.UpdatedAt = now

		if completed && !existingProgress.Completed {
			existingProgress.Completed = true
			existingProgress.CompletedAt = &now
		}

		return s.ProgressRepo.Update(existingProgress)
	}

	// Create new progress
	progress := &models.LearningProgress{
		UserID:       studentID,
		CourseID:     courseID,
		ActivityType: activityType,
		ActivityID:   activityID,
		Score:        &scoreValue,
		MaxScore:     &maxScoreValue,
		Feedback:     feedback,
		GradedBy:     graderId,
		Completed:    completed,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if completed {
		progress.CompletedAt = &now
	}

	return s.ProgressRepo.Create(progress)
}

// GetProgressByID gets a learning progress by ID
func (s *LearningProgressService) GetProgressByID(id uint) (*models.LearningProgress, error) {
	return s.ProgressRepo.FindByID(id)
}

// GetStudentProgress gets all progress for a student in a course
func (s *LearningProgressService) GetStudentProgress(studentID, courseID uint) ([]models.LearningProgress, error) {
	// Verify the student exists
	student, err := s.UserRepo.FindByID(studentID)
	if err != nil {
		return nil, errors.New("student not found")
	}

	// Verify the student is a student
	if student.Role != models.RoleStudent {
		return nil, errors.New("specified user is not a student")
	}

	// Verify the course exists
	_, err = s.CourseRepo.FindByID(courseID)
	if err != nil {
		return nil, errors.New("course not found")
	}

	// Verify the student is enrolled in the course
	_, err = s.EnrollmentRepo.FindByUserAndCourse(studentID, courseID)
	if err != nil {
		return nil, errors.New("student is not enrolled in this course")
	}

	return s.ProgressRepo.GetStudentProgress(studentID, courseID)
}

// GetCourseStudentsProgress gets progress for all students in a course
func (s *LearningProgressService) GetCourseStudentsProgress(courseID uint, requestorID uint) (map[uint][]models.LearningProgress, error) {
	// Verify the requestor exists and is mentor/admin
	requestor, err := s.UserRepo.FindByID(requestorID)
	if err != nil {
		return nil, errors.New("requestor not found")
	}

	if requestor.Role != models.RoleMentor && requestor.Role != models.RoleAdmin {
		return nil, errors.New("only mentors and admins can view all students' progress")
	}

	// Verify the course exists
	course, err := s.CourseRepo.FindByID(courseID)
	if err != nil {
		return nil, errors.New("course not found")
	}

	// For mentors: verify they are assigned to this course
	if requestor.Role == models.RoleMentor && course.MentorID != requestorID {
		return nil, errors.New("mentors can only view progress for courses they are assigned to")
	}

	return s.ProgressRepo.GetCourseStudentsProgress(courseID)
}

// UpdateGradeByMentor allows a mentor to update a grade they previously assigned
func (s *LearningProgressService) UpdateGradeByMentor(
	progressID uint,
	mentorID uint,
	score float64,
	maxScore float64,
	feedback string,
	completed bool,
) error {
	// Verify the mentor exists and is a mentor
	mentor, err := s.UserRepo.FindByID(mentorID)
	if err != nil {
		return errors.New("mentor not found")
	}

	if mentor.Role != models.RoleMentor {
		return errors.New("only mentors can use this function")
	}

	// Get the progress to update
	progress, err := s.ProgressRepo.FindByID(progressID)
	if err != nil {
		return err
	}

	// Verify the course exists
	course, err := s.CourseRepo.FindByID(progress.CourseID)
	if err != nil {
		return errors.New("course not found")
	}

	// Verify mentor is assigned to this course
	if course.MentorID != mentorID {
		return errors.New("mentors can only update grades for their assigned courses")
	}

	// Update the grade
	now := time.Now()
	scoreValue := score
	maxScoreValue := maxScore

	progress.Score = &scoreValue
	progress.MaxScore = &maxScoreValue
	progress.Feedback = feedback
	progress.GradedBy = mentorID
	progress.UpdatedAt = now

	if completed && !progress.Completed {
		progress.Completed = true
		progress.CompletedAt = &now
	}

	return s.ProgressRepo.Update(progress)
}
