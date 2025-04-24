package services

import (
	"LMS/models"
	"LMS/repositories"
	"errors"
	"time"

)

// LearningProgressService menangani logika bisnis kemajuan pembelajaran
type LearningProgressService struct {
	ProgressRepo   *repositories.LearningProgressRepository
	UserRepo       *repositories.UserRepository
	CourseRepo     *repositories.CourseRepository
	EnrollmentRepo *repositories.EnrollmentRepository
	AssignmentRepo *repositories.AssignmentRepository
}

// NewLearningProgressService membuat layanan kemajuan pembelajaran baru
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

// CreateOrUpdateGrade memungkinkan mentor/admin untuk menilai aktivitas siswa
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
	// Verifikasi bahwa grader tersebut ada dan merupakan mentor atau admin
	grader, err := s.UserRepo.FindByID(graderId)
	if err != nil {
		return errors.New("grader not found")
	}

	if grader.Role != models.RoleMentor && grader.Role != models.RoleAdmin {
		return errors.New("only mentors and admins can grade student activities")
	}

	// Verifikasi keberadaan siswa
	student, err := s.UserRepo.FindByID(studentID)
	if err != nil {
		return errors.New("student not found")
	}

	if student.Role != models.RoleStudent {
		return errors.New("only students can be graded")
	}

	// Verifikasi keberadaan kursus
	course, err := s.CourseRepo.FindByID(courseID)
	if err != nil {
		return errors.New("course not found")
	}

	// Untuk mentor: pastikan mereka ditugaskan untuk kursus ini
	if grader.Role == models.RoleMentor && course.MentorID != graderId {
		return errors.New("mentors can only grade courses they are assigned to")
	}
	// Verifikasi bahwa siswa telah terdaftar dalam kursus
	_, err = s.EnrollmentRepo.FindByUserAndCourse(studentID, courseID)
	if err != nil {
		return errors.New("student is not enrolled in this course")
	}

	// Untuk penugasan: verifikasi bahwa penugasan itu ada
	if activityType == models.ProgressTypeAssignment {
		_, err = s.AssignmentRepo.FindByID(activityID)
		if err != nil {
			return errors.New("assignment not found")
		}
	}

	// Periksa apakah kemajuan sudah ada
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
		// Perbarui kemajuan yang ada
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
	// Buat kemajuan baru
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

// GetProgressByID mendapatkan kemajuan pembelajaran berdasarkan ID
func (s *LearningProgressService) GetProgressByID(id uint) (*models.LearningProgress, error) {
	return s.ProgressRepo.FindByID(id)
}

// GetStudentProgress mendapatkan semua kemajuan untuk seorang siswa dalam sebuah kursus
func (s *LearningProgressService) GetStudentProgress(studentID, courseID uint) ([]models.LearningProgress, error) {
	// Verifikasi bahwa siswa adalah seorang siswa
	student, err := s.UserRepo.FindByID(studentID)
	if err != nil {
		return nil, errors.New("student not found")
	}

	// Verifikasi bahwa kursus tersebut ada
	if student.Role != models.RoleStudent {
		return nil, errors.New("specified user is not a student")
	}

	// Verifikasi bahwa kursus tersebut ada
	_, err = s.CourseRepo.FindByID(courseID)
	if err != nil {
		return nil, errors.New("course not found")
	}
	// Verifikasi bahwa siswa telah terdaftar dalam kursus tersebut
	_, err = s.EnrollmentRepo.FindByUserAndCourse(studentID, courseID)
	if err != nil {
		return nil, errors.New("student is not enrolled in this course")
	}

	return s.ProgressRepo.GetStudentProgress(studentID, courseID)
}

// GetCourseStudentsProgress mendapatkan kemajuan untuk semua siswa dalam sebuah kursus
func (s *LearningProgressService) GetCourseStudentsProgress(courseID uint, requestorID uint) (map[uint][]models.LearningProgress, error) {
	// Verifikasi bahwa pemohon ada dan merupakan mentor/admin
	requestor, err := s.UserRepo.FindByID(requestorID)
	if err != nil {
		return nil, errors.New("requestor not found")
	}

	if requestor.Role != models.RoleMentor && requestor.Role != models.RoleAdmin {
		return nil, errors.New("only mentors and admins can view all students' progress")
	}

	// Verifikasi bahwa mata kuliah tersebut ada
	course, err := s.CourseRepo.FindByID(courseID)
	if err != nil {
		return nil, errors.New("course not found")
	}

	// Untuk mentor: verifikasi bahwa mereka ditugaskan untuk mata kuliah ini
	if requestor.Role == models.RoleMentor && course.MentorID != requestorID {
		return nil, errors.New("mentors can only view progress for courses they are assigned to")
	}

	return s.ProgressRepo.GetCourseStudentsProgress(courseID)
}

// UpdateGradeByMentor memungkinkan mentor untuk memperbarui nilai yang mereka berikan sebelumnya
func (s *LearningProgressService) UpdateGradeByMentor(
	progressID uint,
	mentorID uint,
	score float64,
	maxScore float64,
	feedback string,
	completed bool,
) error {
	// Verifikasi bahwa mentor tersebut ada dan merupakan seorang mentor
	mentor, err := s.UserRepo.FindByID(mentorID)
	if err != nil {
		return errors.New("mentor not found")
	}

	if mentor.Role != models.RoleMentor {
		return errors.New("only mentors can use this function")
	}
	// Dapatkan kemajuan untuk diperbarui
	progress, err := s.ProgressRepo.FindByID(progressID)
	if err != nil {
		return err
	}

	// Verifikasi bahwa kursus tersebut ada
	course, err := s.CourseRepo.FindByID(progress.CourseID)
	if err != nil {
		return errors.New("course not found")
	}

	// Verifikasi mentor yang ditugaskan untuk mata kuliah ini
	if course.MentorID != mentorID {
		return errors.New("mentors can only update grades for their assigned courses")
	}

	// Perbarui nilai
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
