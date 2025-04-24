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

// DBConfig represents database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// BuildDBConfig creates database configuration from environment variables
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

// Helper function to get environment variables with default fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// DbURL returns database connection string
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

// InitDB initializes database connection
func InitDB() (*gorm.DB, error) {
	dbConfig := BuildDBConfig()
	dsn := DbURL(dbConfig)

	// Configure logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level (change to logger.Silent in production)
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	// Open connection to database with modified configuration
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                      dsn,
		DisableDatetimePrecision: true, // Disable datetime precision to fix the error
	}), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return nil, err
	}

	// Reset SQL mode to ensure compatibility
	db.Exec("SET @@sql_mode=''")

	// Auto migrate the schema
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
