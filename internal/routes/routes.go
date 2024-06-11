package routes

import (
	"github.com/kamalsingh200238/ubu_management/internal/controllers"
	"github.com/kamalsingh200238/ubu_management/internal/middlewares"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.GET("/", controllers.ShowHomePage)

	e.GET("/login", controllers.ShowLoginPage)
	e.POST("/coordinator-login", controllers.CooridnatorLogin)
	e.POST("/student-login", controllers.StudentLogin)

	coordinatorGroup := e.Group("/coordinator", middlewares.CoordinatorAuthMiddleware)
	coordinatorGroup.GET("", controllers.ShowCoordinatorDashboard)
	coordinatorGroup.GET("/edit-society-modal/:id", controllers.ShowEditSocietyModal)
	coordinatorGroup.PATCH("/edit-society/:id", controllers.EditSociety)
	coordinatorGroup.GET("/create-society-modal", controllers.ShowCreateSocietyModal)
	coordinatorGroup.GET("/create-society", controllers.CreateSociety)
	coordinatorGroup.PATCH("/enable-society/:id", controllers.EnableSociety)
	coordinatorGroup.PATCH("/disable-society/:id", controllers.DisableSociety)

	studentGroup := e.Group("/student", middlewares.StudentAuthMiddleware)
	studentGroup.GET("", controllers.ShowStudentDashboard)
	studentGroup.POST("/enroll-in-society/:studentID/:societyID", controllers.EnrollInSociety)
	studentGroup.POST("/leave-society/:studentID/:societyID", controllers.LeaveSociety)
}
