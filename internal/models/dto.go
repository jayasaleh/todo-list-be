package models

import "time"

// Todo DTOs
type CreateTodoRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	CategoryID  *uint      `json:"category_id"`
	Priority    Priority   `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
}

type UpdateTodoRequest struct {
	Title       *string     `json:"title"`
	Description *string     `json:"description"`
	CategoryID  *uint       `json:"category_id"`
	Priority    *Priority   `json:"priority"`
	Completed   *bool       `json:"completed"`
	DueDate     *time.Time  `json:"due_date"`
}

type TodoResponse struct {
	ID          uint                 `json:"id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Completed   bool                 `json:"completed"`
	Category    *CategoryResponse    `json:"category,omitempty"`
	CategoryID  *uint                `json:"category_id,omitempty"`
	Priority    Priority             `json:"priority"`
	DueDate     *time.Time           `json:"due_date"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

type CategoryResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

// Category DTOs
type CreateCategoryRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color"`
}

type UpdateCategoryRequest struct {
	Name  *string `json:"name"`
	Color *string `json:"color"`
}

// Pagination DTOs
type PaginationParams struct {
	Page      int    `form:"page"`
	Limit     int    `form:"limit"`
	Search    string `form:"search"`
	SortBy    string `form:"sort_by"`
	SortOrder string `form:"sort_order"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
}

// ToTodoResponse converts Todo model to TodoResponse DTO
func ToTodoResponse(todo Todo) TodoResponse {
	response := TodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
		CategoryID:  todo.CategoryID,
		Priority:    todo.Priority,
		DueDate:     todo.DueDate,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}

	if todo.Category != nil {
		response.Category = &CategoryResponse{
			ID:        todo.Category.ID,
			Name:      todo.Category.Name,
			Color:     todo.Category.Color,
			CreatedAt: todo.Category.CreatedAt,
		}
	}

	return response
}

// ToCategoryResponse converts Category model to CategoryResponse DTO
func ToCategoryResponse(category Category) CategoryResponse {
	return CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Color:     category.Color,
		CreatedAt: category.CreatedAt,
	}
}

