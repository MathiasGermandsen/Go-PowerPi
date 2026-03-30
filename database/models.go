package database

import "gorm.io/gorm"

type PowerTable struct {
	gorm.Model
	Price   float64 `json:"price" gorm:"not null"`
	Company string  `json:"company" gorm:"not null"`
}
