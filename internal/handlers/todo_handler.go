package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jayasaleh/todo-list/be/internal/models"
	"github.com/jayasaleh/todo-list/be/internal/services"
	"github.com/jayasaleh/todo-list/be/pkg/utils"
)

type TodoHandler struct {
	todoService *services.TodoService
}

func NewTodoHandler() *TodoHandler {
	return &TodoHandler{
		todoService: services.NewTodoService(),
	}
}

// Get Todos
func (h *TodoHandler) GetTodos(c *gin.Context) {
	var params models.PaginationParams

	if err := c.ShouldBindQuery(&params); err != nil {
		utils.BadRequest(c, "Invalid query parameters")
		return
	}

	todos, pagination, err := h.todoService.GetTodos(params)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	var todoResponses []models.TodoResponse
	for _, todo := range todos {
		todoResponses = append(todoResponses, models.ToTodoResponse(todo))
	}

	if pagination == nil {
		utils.InternalServerError(c, "Pagination data is missing")
		return
	}

	utils.PaginatedResponse(c, "Successfully fetched todos", todoResponses, pagination)
}

// Get Todo by ID
func (h *TodoHandler) GetTodo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid todo ID")
		return
	}

	todo, err := h.todoService.GetTodoByID(uint(id))
	if err != nil {
		if err.Error() == "todo not found" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.OK(c, "Successfully fetching todo", models.ToTodoResponse(*todo))
}

// Create Todo
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var req models.CreateTodoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	todo, err := h.todoService.CreateTodo(req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, "Todo created successfully", models.ToTodoResponse(*todo))
}

// Update Todo
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid todo ID")
		return
	}

	var req models.UpdateTodoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	todo, err := h.todoService.UpdateTodo(uint(id), req)
	if err != nil {
		if err.Error() == "todo not found" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.BadRequest(c, err.Error())
		return
	}

	utils.OK(c, "Todo updated successfully", models.ToTodoResponse(*todo))
}

// Delete Todo
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid todo ID")
		return
	}

	err = h.todoService.DeleteTodo(uint(id))
	if err != nil {
		if err.Error() == "todo not found" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.OK(c, "Todo deleted successfully", nil)
}

// Toggle Todo Complete
func (h *TodoHandler) ToggleComplete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid todo ID")
		return
	}

	todo, err := h.todoService.ToggleComplete(uint(id))
	if err != nil {
		if err.Error() == "todo not found" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.OK(c, "Todo completion status updated successfully", models.ToTodoResponse(*todo))
}
