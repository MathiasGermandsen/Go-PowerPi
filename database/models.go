package database

import (
	"time"

	"gorm.io/gorm"
)

type PowerTable struct {
	gorm.Model
	Price         float64 `json:"price" gorm:"not null"`
	Company       string  `json:"company" gorm:"not null"`
	UserID        string  `json:"userId" gorm:"uniqueIndex;not null"`
	SelectionMode string  `json:"selectionMode" gorm:"not null;default:''"`
	NumberOfHours int     `json:"numberOfHours" gorm:"not null;default:0"`
}

type RevokedToken struct {
	ID        uint      `gorm:"primaryKey"`
	JTI       string    `gorm:"uniqueIndex;not null"`
	RevokedAt time.Time `gorm:"autoCreateTime"`
	ExpiresAt time.Time `gorm:"not null"`
}

type ServiceRegistration struct {
	ID            uint      `gorm:"primaryKey"`
	ServiceName   string    `gorm:"not null"`
	Audience      string    `gorm:"not null"`
	Scopes        string    `gorm:"not null"`
	JTI           string    `gorm:"uniqueIndex;not null"`
	ExpiresAt     time.Time `gorm:"not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	LastRotatedAt time.Time
}
