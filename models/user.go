package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Role type for user roles
type Role string

// User roles
const (
	RoleAdmin   Role = "admin"
	RoleMentor  Role = "mentor"
	RoleStudent Role = "student"
)

// User represents the user entity

type User struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Email     string    `gorm:"size:255;not null;unique" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"` // Password not exported to JSON
	Role      Role      `gorm:"type:enum('admin','mentor','student');not null" json:"role"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"	`
}

// type User struct {
// 	gorm.Model
// 	ID        uint           `gorm:"primaryKey" json:"id"`
// 	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
// 	Email     string         `gorm:"type:varchar(255);unique;not null" json:"email"`
// 	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
// 	Role      string         `gorm:"type:enum('admin','mentor','student');not null" json:"role"`
// 	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
// 	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
// 	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
// }

// BeforeSave hashes the user password before saving to database
func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// CheckPassword checks if the provided password matches the stored hash
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
