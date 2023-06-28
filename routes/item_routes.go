package routes

import (
	"github.com/aldoramirezmartinez/fiber-api/controllers"
	"github.com/gofiber/fiber/v2"
)

type ItemRoutes struct {
	router         fiber.Router
	itemController *controllers.ItemController
}

func NewItemRoutes(router fiber.Router, itemController *controllers.ItemController) *ItemRoutes {
	return &ItemRoutes{
		router:         router,
		itemController: itemController,
	}
}

func (ir *ItemRoutes) SetupRoutes() {
	itemRouter := ir.router.Group("/api/items")

	itemRouter.Get("/", ir.itemController.GetAllItems)
	itemRouter.Get("/:id", ir.itemController.GetItem)
	itemRouter.Post("/", ir.itemController.CreateItem)
	itemRouter.Put("/:id", ir.itemController.UpdateItem)
	itemRouter.Delete("/:id", ir.itemController.DeleteItem)

}
