package config

import (
	"LMS/models"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBConfig konfigurasi basis data
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// BuildDBConfig membuat konfigurasi basis data dari variabel lingkungan
func BuildDBConfig() *DBConfig {
	dbConfig := DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "Yosua"),
		DBName:   getEnv("DB_NAME", "lmsba"),
	}
	return &dbConfig
}

// Fungsi pembantu untuk mendapatkan variabel lingkungan dengan fallback default
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// DbURL mengembalikan string koneksi database
func DbURL(dbConfig *DBConfig) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)
}

// InitDB menginisialisasi koneksi database
func InitDB() (*gorm.DB, error) {
	dbConfig := BuildDBConfig()
	dsn := DbURL(dbConfig)

	// Konfigurasikan  logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), //mencatat io
		logger.Config{
			SlowThreshold:             time.Second, // Ambang batas SQL
			LogLevel:                  logger.Info, // Level log (ubah ke logger.Silent dalam produksi)
			IgnoreRecordNotFoundError: true,        // Abaikan kesalahan ErrRecordNotFound untuk pencatat
			Colorful:                  true,        // Nonaktifkan warna
		},
	)

	//Buka koneksi ke database dengan konfigurasi yang dimodifikasi
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                      dsn,
		DisableDatetimePrecision: true, // Nonaktifkan presisi waktu data untuk memperbaiki kesalahan
	}), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return nil, err
	}

	//Atur ulang mode SQL untuk memastikan kompatibilitas
	db.Exec("SET @@sql_mode=''")

	// migrasi basis data
	// Menggunakan AutoMigrate untuk membuat tabel dan kolom yang diperlukan
	err = db.AutoMigrate(
		&models.User{},
		&models.Course{},
		&models.Enrollment{},
		&models.Material{},
		&models.Assignment{},
		&models.Submission{},
		&models.Assessment{},
		&models.Discussion{},
		&models.Comment{},
		&models.LearningProgress{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
