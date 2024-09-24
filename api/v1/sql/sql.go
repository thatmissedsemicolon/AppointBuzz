package sql

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
    ID        uuid.UUID  `gorm:"type:uuid;primary_key;" json:"id"`
    Email     string     `gorm:"type:varchar(100);unique_index" json:"email"`
    Password  string     `gorm:"size:60" json:"-"`
    CreatedAt time.Time  `gorm:"type:timestamp" json:"created_at"`
    UpdatedAt time.Time  `gorm:"type:timestamp" json:"updated_at"`
    DeletedAt *time.Time `gorm:"type:timestamp" json:"deleted_at,omitempty"`
}
