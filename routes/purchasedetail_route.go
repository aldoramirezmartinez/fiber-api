package routes

import (
	"github.com/aldoramirezmartinez/fiber-api/controllers"
	"github.com/gofiber/fiber/v2"
)

type PurchaseDetailRoutes struct {
	router                   fiber.Router
	purchaseDetailController *controllers.PurchaseDetailController
}

func NewPurchaseDetailRoutes(router fiber.Router, purchaseDetailController *controllers.PurchaseDetailController) *PurchaseDetailRoutes {
	return &PurchaseDetailRoutes{
		router:                   router,
		purchaseDetailController: purchaseDetailController,
	}
}

func (pdr *PurchaseDetailRoutes) SetupRoutes() {
	purchaseDetailRouter := pdr.router.Group("/api/purchasedetails")

	purchaseDetailRouter.Get("/", pdr.purchaseDetailController.GetAllPurchaseDetails)
	purchaseDetailRouter.Get("/:id", pdr.purchaseDetailController.GetPurchaseDetail)
	purchaseDetailRouter.Post("/", pdr.purchaseDetailController.CreatePurchaseDetail)
	purchaseDetailRouter.Put("/:id", pdr.purchaseDetailController.UpdatePurchaseDetail)
	purchaseDetailRouter.Delete("/:id", pdr.purchaseDetailController.DeletePurchaseDetail)
}
