package services

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/jayasaleh/todo-list/be/internal/database"
	"github.com/jayasaleh/todo-list/be/internal/models"
)

type TodoService struct {
	db *gorm.DB
}

func NewTodoService() *TodoService {
	return &TodoService{
		db: database.GetDB(),
	}
}

// Create Todo
func (s *TodoService) CreateTodo(req models.CreateTodoRequest) (*models.Todo, error) {
	todo := models.Todo{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	}

	// Validate category exists
	var category models.Category
	if err := s.db.First(&category, req.CategoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, fmt.Errorf("failed to validate category: %w", err)
	}
	categoryID := req.CategoryID
	todo.CategoryID = &categoryID

	// Validate priority (required)
	if !models.ValidatePriority(req.Priority) {
		return nil, errors.New("invalid priority value. Must be 'high', 'medium', or 'low'")
	}
	todo.Priority = req.Priority

	if err := s.db.Create(&todo).Error; err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	// Preload category since category_id is required
	s.db.Preload("Category").First(&todo, todo.ID)

	return &todo, nil
}

// Get Todo by ID
func (s *TodoService) GetTodoByID(id uint) (*models.Todo, error) {
	var todo models.Todo

	if err := s.db.Preload("Category").First(&todo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("todo not found")
		}
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	return &todo, nil
}

// Get Todos with Pagination
func (s *TodoService) GetTodos(params models.PaginationParams) ([]models.Todo, *models.Pagination, error) {
	var todos []models.Todo
	var total int64

	query := s.db.Model(&models.Todo{}).Preload("Category")

	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("LOWER(title) LIKE LOWER(?)", searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to count todos: %w", err)
	}

	page := params.Page
	if page < 1 {
		page = 1
	}

	limit := params.Limit
	if limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	offset := (page - 1) * limit

	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}

	sortOrder := params.SortOrder
	if sortOrder == "" {
		sortOrder = "desc"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))

	if err := query.Offset(offset).Limit(limit).Find(&todos).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to get todos: %w", err)
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	pagination := &models.Pagination{
		CurrentPage: page,
		PerPage:     limit,
		Total:       total,
		TotalPages:  totalPages,
	}

	return todos, pagination, nil
}

// Update Todo
func (s *TodoService) UpdateTodo(id uint, req models.UpdateTodoRequest) (*models.Todo, error) {
	var todo models.Todo

	if err := s.db.First(&todo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("todo not found")
		}
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}
	if req.CategoryID != nil {
		var category models.Category
		if err := s.db.First(&category, *req.CategoryID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("category not found")
			}
			return nil, fmt.Errorf("failed to validate category: %w", err)
		}
		todo.CategoryID = req.CategoryID
	}
	if req.Priority != nil {
		if !models.ValidatePriority(*req.Priority) {
			return nil, errors.New("invalid priority value. Must be 'high', 'medium', or 'low'")
		}
		todo.Priority = *req.Priority
	}
	if req.DueDate != nil {
		todo.DueDate = req.DueDate
	}

	if err := s.db.Save(&todo).Error; err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	if todo.CategoryID != nil {
		s.db.Preload("Category").First(&todo, todo.ID)
	}

	return &todo, nil
}

// DeleteTodo - Delete todo
func (s *TodoService) DeleteTodo(id uint) error {
	var todo models.Todo

	// Find by ID
	if err := s.db.First(&todo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("todo not found")
		}
		return fmt.Errorf("failed to get todo: %w", err)
	}

	// Delete from DB
	if err := s.db.Delete(&todo).Error; err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	return nil
}

// Toggle Todo Complete
func (s *TodoService) ToggleComplete(id uint) (*models.Todo, error) {
	var todo models.Todo

	if err := s.db.First(&todo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("todo not found")
		}
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	todo.Completed = !todo.Completed

	if err := s.db.Save(&todo).Error; err != nil {
		return nil, fmt.Errorf("failed to toggle todo completion: %w", err)
	}

	if todo.CategoryID != nil {
		s.db.Preload("Category").First(&todo, todo.ID)
	}

	return &todo, nil
}
