package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

)

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleMentor  Role = "mentor"
	RoleStudent Role = "student"
)

type User struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Email     string    `gorm:"size:255;not null;unique" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Role      Role      `gorm:"type:enum('admin','mentor','student');not null" json:"role"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"	`
}

// BeforeSave melakuakan penyimpanan/ memasukan sandi pengguna sebelum menyimpan database
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

// CheckPassword memeriksa apakah kata sandi cocok dengan  yang disimpan
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
