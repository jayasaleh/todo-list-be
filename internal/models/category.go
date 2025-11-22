package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null;uniqueIndex"`
	Color     string         `json:"color" gorm:"default:'#3B82F6'"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationship
	Todos []Todo `json:"-" gorm:"foreignKey:CategoryID"`
}

// TableName specifies the table name for Category model
func (Category) TableName() string {
	return "categories"
}
