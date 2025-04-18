package config

import (
	"LMS/models"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB membuka koneksi ke database MySQL
func ConnectDB() *gorm.DB {
	// Ganti username, password, host, dan database sesuai konfigurasi Anda
	dsn := "root:Yosua@tcp(127.0.0.1:3306)/lmsbar?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terkoneksi ke database:", err)
	}
	DB = db
	return db
}

// MigrateDB melakukan migrasi semua model
func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Course{},
		&models.Enrollment{},
		&models.Material{},
		&models.Assignment{},
		&models.Submission{},
		&models.Assessment{},
		&models.Discussion{},
		&models.Comment{},
	)
	if err != nil {
		fmt.Println("Migrasi gagal:", err)
	} else {
		fmt.Println("Migrasi berhasil")
	}
}
