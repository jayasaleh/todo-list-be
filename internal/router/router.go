package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jayasaleh/todo-list/be/internal/handlers"
	"github.com/jayasaleh/todo-list/be/internal/middleware"
	"github.com/jayasaleh/todo-list/be/pkg/utils"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())

	router.GET("/health", func(c *gin.Context) {
		utils.OK(c, "Todo List API is running", nil)
	})

	todoHandler := handlers.NewTodoHandler()
	categoryHandler := handlers.NewCategoryHandler()

	api := router.Group("/api")
	{
		todos := api.Group("/todos")
		{
			todos.GET("", todoHandler.GetTodos)
			todos.GET("/:id", todoHandler.GetTodo)
			todos.POST("", todoHandler.CreateTodo)
			todos.PUT("/:id", todoHandler.UpdateTodo)
			todos.DELETE("/:id", todoHandler.DeleteTodo)
			todos.PATCH("/:id/complete", todoHandler.ToggleComplete)
		}

		categories := api.Group("/categories")
		{
			categories.GET("", categoryHandler.GetCategories)
			categories.GET("/:id", categoryHandler.GetCategory)
			categories.POST("", categoryHandler.CreateCategory)
			categories.PUT("/:id", categoryHandler.UpdateCategory)
			categories.DELETE("/:id", categoryHandler.DeleteCategory)
		}
	}

	return router
}
