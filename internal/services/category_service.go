package services

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/jayasaleh/todo-list/be/internal/database"
	"github.com/jayasaleh/todo-list/be/internal/models"
)

type CategoryService struct {
	db *gorm.DB
}

func NewCategoryService() *CategoryService {
	return &CategoryService{
		db: database.GetDB(),
	}
}

// Create Category
func (s *CategoryService) CreateCategory(req models.CreateCategoryRequest) (*models.Category, error) {
	color := req.Color
	if color == "" {
		color = "#3B82F6"
	}

	category := models.Category{
		Name:  req.Name,
		Color: color,
	}

	if err := s.db.Create(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errors.New("category name already exists")
		}
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return &category, nil
}

// Get Category by ID
func (s *CategoryService) GetCategoryByID(id uint) (*models.Category, error) {
	var category models.Category

	if err := s.db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return &category, nil
}

// Get All Categories
func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category

	if err := s.db.Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	return categories, nil
}

// Update Category
func (s *CategoryService) UpdateCategory(id uint, req models.UpdateCategoryRequest) (*models.Category, error) {
	var category models.Category

	if err := s.db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Color != nil {
		category.Color = *req.Color
	}

	if err := s.db.Save(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errors.New("category name already exists")
		}
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return &category, nil
}

// Delete Category
func (s *CategoryService) DeleteCategory(id uint) error {
	var category models.Category

	if err := s.db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return fmt.Errorf("failed to get category: %w", err)
	}

	var count int64
	if err := s.db.Model(&models.Todo{}).Where("category_id = ?", id).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check category usage: %w", err)
	}

	if count > 0 {
		return errors.New("cannot delete category that is being used by todos")
	}

	if err := s.db.Delete(&category).Error; err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}
