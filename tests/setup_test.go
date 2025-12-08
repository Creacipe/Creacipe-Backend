package tests

import (
	"creacipe-backend/config"
	"creacipe-backend/models"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"

	// Use pure Go SQLite driver (no CGO required)
	"github.com/glebarez/sqlite"
)

// SetupTestDB menyiapkan database SQLite di RAM dengan SEMUA TABEL
func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("Gagal connect ke database test: " + err.Error())
	}

	// Override variabel global di package config
	config.DB = db

	// Manual create tables dengan raw SQL untuk bypass enum issue
	// SQLite akan gunakan TEXT untuk semua enum fields
	
	db.Exec(`CREATE TABLE IF NOT EXISTS roles (
		role_id INTEGER PRIMARY KEY,
		role_name VARCHAR(255) NOT NULL UNIQUE
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS users (
		user_id INTEGER PRIMARY KEY AUTOINCREMENT,
		role_id INTEGER NOT NULL DEFAULT 3,
		status_user VARCHAR(20) DEFAULT 'active',
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		created_at DATETIME,
		updated_at DATETIME,
		FOREIGN KEY (role_id) REFERENCES roles(role_id)
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS user_profiles (
		profile_id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL UNIQUE,
		bio TEXT,
		phone_number VARCHAR(20),
		address TEXT,
		profile_picture VARCHAR(255),
		profile_picture_url VARCHAR(255),
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS categories (
		category_id INTEGER PRIMARY KEY AUTOINCREMENT,
		category_name VARCHAR(100) NOT NULL UNIQUE,
		created_at DATETIME,
		updated_at DATETIME
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS tags (
		tag_id INTEGER PRIMARY KEY AUTOINCREMENT,
		category_id INTEGER NOT NULL,
		tag_name VARCHAR(100) NOT NULL UNIQUE,
		created_at DATETIME,
		updated_at DATETIME,
		FOREIGN KEY (category_id) REFERENCES categories(category_id)
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS menu_tags (
		menu_id INTEGER NOT NULL,
		tag_id INTEGER NOT NULL,
		PRIMARY KEY (menu_id, tag_id),
		FOREIGN KEY (menu_id) REFERENCES menus(menu_id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags(tag_id) ON DELETE CASCADE
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS menus (
		menu_id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		ingredients TEXT,
		instructions TEXT,
		image_url VARCHAR(255),
		status VARCHAR(20) DEFAULT 'pending',
		rejection_reason TEXT,
		created_at DATETIME,
		updated_at DATETIME,
		FOREIGN KEY (user_id) REFERENCES users(user_id)
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS menu_votes (
		vote_id INTEGER PRIMARY KEY AUTOINCREMENT,
		menu_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		likes_count INTEGER DEFAULT 0 NOT NULL,
		dislikes_count INTEGER DEFAULT 0 NOT NULL,
		created_at DATETIME,
		updated_at DATETIME,
		FOREIGN KEY (menu_id) REFERENCES menus(menu_id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS comments (
		comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
		menu_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		parent_id INTEGER,
		comment_text TEXT NOT NULL,
		created_at DATETIME,
		FOREIGN KEY (menu_id) REFERENCES menus(menu_id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
		FOREIGN KEY (parent_id) REFERENCES comments(comment_id) ON DELETE CASCADE
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS notifications (
		notification_id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title VARCHAR(255) NOT NULL,
		message TEXT NOT NULL,
		type VARCHAR(20) DEFAULT 'info',
		is_read BOOLEAN DEFAULT 0,
		related_id INTEGER,
		related_type VARCHAR(50),
		created_at DATETIME,
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS log_activities (
		log_id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		action VARCHAR(255) NOT NULL,
		related_id INTEGER,
		related_table VARCHAR(50),
		created_at DATETIME,
		FOREIGN KEY (user_id) REFERENCES users(user_id)
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS password_resets (
		reset_id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		verification_code VARCHAR(6) NOT NULL,
		expires_at DATETIME NOT NULL,
		is_used BOOLEAN DEFAULT 0,
		created_at DATETIME,
		FOREIGN KEY (user_id) REFERENCES users(user_id)
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS user_bookmarks (
		user_id INTEGER NOT NULL,
		menu_id INTEGER NOT NULL,
		PRIMARY KEY (user_id, menu_id),
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
		FOREIGN KEY (menu_id) REFERENCES menus(menu_id) ON DELETE CASCADE
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS email_verifications (
		verification_id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		new_email VARCHAR(255) NOT NULL,
		verification_code VARCHAR(6) NOT NULL,
		expires_at DATETIME NOT NULL,
		is_used BOOLEAN DEFAULT 0,
		created_at DATETIME,
		FOREIGN KEY (user_id) REFERENCES users(user_id)
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS notifications (
		notification_id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title VARCHAR(255) NOT NULL,
		message TEXT NOT NULL,
		type VARCHAR(20) DEFAULT 'info',
		is_read BOOLEAN DEFAULT 0,
		related_id INTEGER,
		related_type VARCHAR(50),
		created_at DATETIME,
		FOREIGN KEY (user_id) REFERENCES users(user_id)
	)`)

	// Seed default roles untuk testing
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.Role{RoleID: 1, RoleName: "admin"})
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.Role{RoleID: 2, RoleName: "editor"})
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.Role{RoleID: 3, RoleName: "member"})

	return db
}

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func TestMain(m *testing.M) {
	// Set Environment Variables untuk testing
	os.Setenv("RECAPTCHA_SECRET_KEY", "dummy-secret")
	os.Setenv("JWT_SECRET", "rahasia-test-123")
	os.Setenv("SMTP_HOST", "smtp.test.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_USERNAME", "test@example.com")
	os.Setenv("SMTP_PASSWORD", "dummy-password")
	os.Setenv("SMTP_FROM_EMAIL", "noreply@test.com")
	os.Setenv("SMTP_FROM_NAME", "Test App")

	exitVal := m.Run()
	os.Exit(exitVal)
}