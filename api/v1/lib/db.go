package lib

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitializeDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := DB.AutoMigrate(&User{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}

type User struct {
	ID             uuid.UUID      `gorm:"type:text;primaryKey;" json:"id"`
	Name           string         `gorm:"type:varchar(100);not null" json:"name"`
	Email          string         `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password       string         `gorm:"type:text;not null" json:"-"`
	Roles          string         `gorm:"type:text;not null;default:''" json:"-"`
	Status         string         `gorm:"type:text;not null;default:'active'" json:"-"`
	ProfilePicture *string        `gorm:"type:text" json:"profile_picture,omitempty"`
	LastLoginAt    *time.Time     `gorm:"type:timestamp" json:"-"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"-"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"-"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
