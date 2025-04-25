# E-Learning Service API

A comprehensive Learning Management System (LMS) built with Go (Golang), featuring user management, course creation, material handling, assignments, submissions, assessments, discussions, and learning progress tracking.

---

## 🚀 Technology Stack

- **Backend:** Gin (Go)
- **ORM:** GORM
- **Database:** MySQL
- **Authentication:** JWT
- **Containerization:** Docker & Docker Compose
- **Go Version:** 1.20+

---

## ✨ Features

- **User Management:**
  - Registration, login, and role-based access control (Admin, Mentor, Student).
- **Course Management:**
  - CRUD operations on courses.
- **Learning Materials:**
  - Upload, download, and manage PDFs/videos.
- **Assignments & Submissions:**
  - Create assignments with due dates and grading criteria; students submit work.
- **Assessments:**
  - Mentors grade submissions and provide feedback.
- **Discussions & Comments:**
  - Forum-like threads per course and nested comments.
- **Learning Progress:**
  - Track student activity (material access, submissions) and generate progress reports.

---

## 📁 Project Structure

```
LMS/
├── cmd/
│   └── main.go                # Entry point
├── config/
│   └── database.go            # DB connection
├── controllers/               # HTTP handlers
│   ├── auth_controller.go
│   ├── course_controller.go
│   └── ...
├── middlewares/
│   └── jwt.go                 # JWT auth
├── models/                    # GORM models
│   └── user.go, course.go, ...
├── repositories/              # Data access layer
│   └── user_repository.go, ...
├── services/                  # Business logic layer
│   └── auth_service.go, ...
├── routes/
│   └── routes.go              # Route definitions
├── utils/
│   └── jwt.go                 # Token utils
├── uploads/
│   ├── materials/
│   └── submissions/
├── Dockerfile
├── docker-compose.yml
├── init.sql                   # DB schema
├── go.mod, go.sum
└── postman_collection.json    # API tests
```

---

## 🗄️ Database Schema

- **One-to-Many:**
  - `Users → Courses` (Mentor can create many)
  - `Courses → Materials, Assignments, Discussions`
  - `Discussions → Comments`
- **Many-to-Many:**
  - `Users ↔ Courses` via `Enrollments`
- **One-to-One:**
  - `Submissions ↔ Assessments`

---

## 🔄 Business Flow

1. **Admin:** Creates courses and invites mentors/students.
2. **Mentor:** Uploads materials, creates assignments.
3. **Student:** Enrolls in courses, accesses materials, submits assignments.
4. **System:** Tracks activities and generates progress reports.

---

## 📬 API Endpoints

> **Full documentation** in `postman_collection.json` or inspect `controllers/` directory.

### Authentication

| Method | Endpoint                    | Description            |
| ------ | --------------------------- | ---------------------- |
| POST   | `/api/auth/register`        | Register new user      |
| POST   | `/api/auth/login`           | User login             |
| GET    | `/api/profile`              | Get current profile    |

### Courses

| Method | Endpoint                            | Description                 |
| ------ | ----------------------------------- | --------------------------- |
| POST   | `/api/courses`                      | Create course               |
| GET    | `/api/courses`                      | List all courses            |
| GET    | `/api/courses/:id`                  | Get course by ID            |
| GET    | `/api/courses/mentor/:mentorId`     | Get courses by mentor       |
| PUT    | `/api/courses/:id`                  | Update course               |
| DELETE | `/api/courses/:id`                  | Delete course               |

#### Enrollments

| Method | Endpoint                                                   | Description                        |
| ------ | ---------------------------------------------------------- | ---------------------------------- |
| POST   | `/api/enrollments`                                         | Create a new enrollment           |
| GET    | `/api/enrollments/check?user_id={userId}&course_id={courseId}` | Check if a user is enrolled       |
| GET    | `/api/enrollments/{id}`                                    | Retrieve enrollment by ID         |
| GET    | `/api/enrollments/user/{userId}`                           | List all enrollments for a user   |
| GET    | `/api/enrollments/course/{courseId}`                       | List all enrollments for a course |
| DELETE | `/api/enrollments/{id}`                                    | Delete an enrollment              |

#### Materials

| Method | Endpoint                                 | Description                           |
| ------ | ---------------------------------------- | ------------------------------------- |
| POST   | `/api/materials`                         | Upload a new material                |
| GET    | `/api/materials/{id}`                    | Retrieve material by ID              |
| GET    | `/api/materials/course/{courseId}`       | List all materials for a course      |
| GET    | `/api/materials/download/{id}`           | Download material file               |
| PUT    | `/api/materials/{id}`                    | Update material metadata             |
| DELETE | `/api/materials/{id}`                    | Delete a material                    |

#### Assignments

| Method | Endpoint                                 | Description                           |
| ------ | ---------------------------------------- | ------------------------------------- |
| POST   | `/api/assignments`                       | Create a new assignment              |
| GET    | `/api/assignments/{id}`                  | Retrieve assignment by ID            |
| GET    | `/api/assignments/course/{courseId}`     | List all assignments for a course    |
| PUT    | `/api/assignments/{id}`                  | Update an assignment                 |
| DELETE | `/api/assignments/{id}`                  | Delete an assignment                 |

#### Submissions

| Method | Endpoint                                              | Description                              |
| ------ | ----------------------------------------------------- | ---------------------------------------- |
| POST   | `/api/submissions`                                    | Submit assignment work                  |
| GET    | `/api/submissions/{id}`                               | Retrieve submission by ID               |
| GET    | `/api/submissions/assignment/{assignmentId}`          | List all submissions for an assignment  |
| GET    | `/api/submissions/student/{studentId}`                | List all submissions by a student       |
| GET    | `/api/submissions/download/{id}`                      | Download submission file                |
| DELETE | `/api/submissions/{id}`                               | Delete a submission                     |

#### Assessments

| Method | Endpoint                                          | Description                          |
| ------ | ------------------------------------------------- | ------------------------------------ |
| POST   | `/api/assessments`                                | Create a new assessment              |
| GET    | `/api/assessments/{id}`                           | Retrieve assessment by ID            |
| GET    | `/api/assessments/submission/{submissionId}`      | Get assessment for a submission      |
| PUT    | `/api/assessments/{id}`                           | Update an assessment                 |
| DELETE | `/api/assessments/{id}`                           | Delete an assessment                 |

#### Discussions

| Method | Endpoint                                   | Description                           |
| ------ | ------------------------------------------ | ------------------------------------- |
| POST   | `/api/discussions`                         | Start a new discussion thread        |
| GET    | `/api/discussions/course/{courseId}`       | List discussions for a course        |
| GET    | `/api/discussions/{id}`                    | Retrieve discussion by ID            |
| PUT    | `/api/discussions/{id}`                    | Update a discussion                  |
| DELETE | `/api/discussions/{id}`                    | Delete a discussion                  |

#### Comments

| Method | Endpoint                                      | Description                      |
| ------ | --------------------------------------------- | -------------------------------- |
| POST   | `/api/comments`                               | Add a comment to a discussion    |
| GET    | `/api/comments/discussion/{discussionId}`     | List comments for a discussion   |
| GET    | `/api/comments/{id}`                          | Retrieve comment by ID           |
| PUT    | `/api/comments/{id}`                          | Update a comment                 |
| DELETE | `/api/comments/{id}`                          | Delete a comment                 |

#### Progress

| Method | Endpoint                                              | Description                                  |
| ------ | ----------------------------------------------------- | -------------------------------------------- |
| POST   | `/api/progress/grade`                                 | Record student activity and score            |
| PUT    | `/api/progress/update-grade/{progressId}`             | Update an existing progress record           |
| GET    | `/api/progress/course/{courseId}/students`            | Get all students' progress in a course       |
| GET    | `/api/progress/student/{studentId}/course/{courseId}` | Get a specific student's progress in course  |
| GET    | `/api/progress/{id}`                                  | Retrieve progress by ID                     |

---


## 🐳 Running with Docker

1. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/lms.git
   cd lms
   ```
2. **Build and start containers**:
   ```bash
   docker-compose up -d
   ```
3. **Access the API** at `http://localhost:8081`

## 🛠️ Manual Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/lms.git
   cd lms
   ```
2. **Configure the database**:
   - Create a MySQL database named `lms`
   - Import the schema:
     ```bash
     mysql -u root -p lms < init.sql
     ```
   - Update connection details in `config/database.go` or via environment variables.
3. **Build and run the application**:
   ```bash
   go build -o lms ./cmd/main.go
   ./lms
   ```
4. **Access the API** at `http://localhost:8080`

## ✅ Testing

- Import `postman_collection.json` into Postman.
- Run **Authentication** tests first to obtain a JWT token.
- Run other endpoint tests; they will use the token automatically.

## 🔒 Security

- Passwords are hashed with **bcrypt**.
- **JWT** tokens are used for authentication.
- Role-based access control enforced on all endpoints.
- File uploads restricted to `uploads/` directories.
- Input validation applied on all requests.

## 📄 License

This project is licensed under the **MIT License**. See [LICENSE](LICENSE) for details.