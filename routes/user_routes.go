package routes

import (
	"github.com/aldoramirezmartinez/fiber-api/controllers"
	"github.com/gofiber/fiber/v2"
)

type UserRoutes struct {
	router         fiber.Router
	userController *controllers.UserController
}

func NewUserRoutes(router fiber.Router, userController *controllers.UserController) *UserRoutes {
	return &UserRoutes{
		router:         router,
		userController: userController,
	}
}

func (ur *UserRoutes) SetupRoutes() {
	userRouter := ur.router.Group("/api/users")

	userRouter.Get("/", ur.userController.GetAllUsers)
	userRouter.Get("/:id", ur.userController.GetUser)
	userRouter.Post("/", ur.userController.CreateUser)
	userRouter.Put("/:id", ur.userController.UpdateUser)
	userRouter.Delete("/:id", ur.userController.DeleteUser)
}
