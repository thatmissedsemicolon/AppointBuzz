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
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := DB.AutoMigrate(&User{}, &FormConfig{}, &Form{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database initialized and migrated successfully!")
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

type FormConfig struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;" json:"id"`
	UserEmail      string         `gorm:"type:varchar(100);not null;index;" json:"user_email"`
	AllowedDomains string         `gorm:"type:text;not null" json:"allowed_domains"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"-"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"-"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type Form struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	FormConfigID uuid.UUID `gorm:"type:uuid;not null;index;" json:"form_id"`
	Data         string    `gorm:"type:text;not null" json:"data"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

func (f *FormConfig) BeforeCreate(tx *gorm.DB) (err error) {
	f.ID = uuid.New()
	return
}
