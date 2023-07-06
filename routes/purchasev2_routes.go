package routes

import (
	"github.com/aldoramirezmartinez/fiber-api/controllers"
	"github.com/gofiber/fiber/v2"
)

type PurchaseV2Routes struct {
	router               fiber.Router
	PurchaseV2Controller *controllers.PurchaseV2Controller
}

func NewPurchaseV2Routes(router fiber.Router, purchaseV2Controller *controllers.PurchaseV2Controller) *PurchaseV2Routes {
	return &PurchaseV2Routes{
		router:               router,
		PurchaseV2Controller: purchaseV2Controller,
	}
}

func (pr *PurchaseV2Routes) SetupRoutes() {
	purchasev2Router := pr.router.Group("/api/purchases")

	purchasev2Router.Get("/", pr.PurchaseV2Controller.GetAllPurchasesV2)
	purchasev2Router.Get("/:purchase_order", pr.PurchaseV2Controller.GetPurchaseV2)
	purchasev2Router.Post("/", pr.PurchaseV2Controller.CreatePurchaseV2)
}
