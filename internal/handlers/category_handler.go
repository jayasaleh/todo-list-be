package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jayasaleh/todo-list/be/internal/models"
	"github.com/jayasaleh/todo-list/be/internal/services"
	"github.com/jayasaleh/todo-list/be/pkg/utils"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		categoryService: services.NewCategoryService(),
	}
}

// Get Categories - GET /api/categories
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	var categoryResponses []models.CategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, models.ToCategoryResponse(category))
	}

	utils.OK(c, "Successfully fetching categories", categoryResponses)
}

// Get Category by ID - GET /api/categories/:id
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid category ID")
		return
	}

	category, err := h.categoryService.GetCategoryByID(uint(id))
	if err != nil {
		if err.Error() == "category not found" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.OK(c, "Successfully fetching category", models.ToCategoryResponse(*category))
}

// Create Category - POST /api/categories
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req models.CreateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	category, err := h.categoryService.CreateCategory(req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, "Category created successfully", models.ToCategoryResponse(*category))
}

// Update Category - PUT /api/categories/:id
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid category ID")
		return
	}

	var req models.UpdateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	category, err := h.categoryService.UpdateCategory(uint(id), req)
	if err != nil {
		if err.Error() == "category not found" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.BadRequest(c, err.Error())
		return
	}

	utils.OK(c, "Category updated successfully", models.ToCategoryResponse(*category))
}

// Delete Category - DELETE /api/categories/:id
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid category ID")
		return
	}

	err = h.categoryService.DeleteCategory(uint(id))
	if err != nil {
		if err.Error() == "category not found" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.BadRequest(c, err.Error())
		return
	}

	utils.OK(c, "Category deleted successfully", nil)
}
