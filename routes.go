package routes

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/yourusername/user-management-app/backend/controllers"
)

func RegisterRoutes(e *echo.Echo, uc *controllers.UserController) {
	e.GET("/users", uc.GetUsers)
	e.POST("/users", uc.CreateUser)
	e.PUT("/users/:id", uc.UpdateUser)
	e.DELETE("/users/:id", uc.DeleteUser)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
