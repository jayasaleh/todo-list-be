package models

import (
	"time"

	"gorm.io/gorm"
)

type Priority string

const (
	PriorityHigh   Priority = "high"
	PriorityMedium Priority = "medium"
	PriorityLow    Priority = "low"
)

type Todo struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	Completed   bool           `json:"completed" gorm:"default:false"`
	CategoryID  *uint          `json:"category_id"`
	Priority    Priority       `json:"priority" gorm:"default:'medium'"`
	DueDate     *time.Time     `json:"due_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationship
	Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

// TableName specifies the table name for Todo model
func (Todo) TableName() string {
	return "todos"
}

// ValidatePriority validates if the priority value is valid
func ValidatePriority(p Priority) bool {
	return p == PriorityHigh || p == PriorityMedium || p == PriorityLow
}
