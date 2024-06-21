package user

import (
	"time"
)

type User struct {
    ID             int       `gorm:"column:id;primaryKey;autoIncrement"`
    Name           string    `gorm:"column:name"`
    Email          string    `gorm:"column:email"`
    PasswordHash   string    `gorm:"column:password_hash"`
    AvatarFileName string    `gorm:"column:avatar_file_name"`
    Role           string    `gorm:"column:role"`
    Token          string    `gorm:"column:token"`
    CreatedAt      time.Time `gorm:"column:created_at"`
    UpdatedAt      time.Time `gorm:"column:updated_at"`
    DeletedAt      *time.Time `gorm:"column:deleted_at"` 
}