package routes

import (
	"github.com/aldoramirezmartinez/fiber-api/controllers"
	"github.com/gofiber/fiber/v2"
)

type PurchaseRoutes struct {
	router             fiber.Router
	purchaseController *controllers.PurchaseController
}

func NewPurchaseRoutes(router fiber.Router, purchaseController *controllers.PurchaseController) *PurchaseRoutes {
	return &PurchaseRoutes{
		router:             router,
		purchaseController: purchaseController,
	}
}

func (pr *PurchaseRoutes) SetupRoutes() {
	purchaseRouter := pr.router.Group("/api/purchases")

	purchaseRouter.Get("/", pr.purchaseController.GetAllPurchases)
	purchaseRouter.Get("/:id", pr.purchaseController.GetPurchase)
	purchaseRouter.Post("/", pr.purchaseController.CreatePurchase)
	purchaseRouter.Put("/:id", pr.purchaseController.UpdatePurchase)
	purchaseRouter.Delete("/:id", pr.purchaseController.DeletePurchase)
}
