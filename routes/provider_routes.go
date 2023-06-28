package routes

import (
	"github.com/aldoramirezmartinez/fiber-api/controllers"
	"github.com/gofiber/fiber/v2"
)

type ProviderRoutes struct {
	router             fiber.Router
	providerController *controllers.ProviderController
}

func NewProviderRoutes(router fiber.Router, providerController *controllers.ProviderController) *ProviderRoutes {
	return &ProviderRoutes{
		router:             router,
		providerController: providerController,
	}
}

func (pr *ProviderRoutes) SetupRoutes() {
	providerRouter := pr.router.Group("/api/providers")

	providerRouter.Get("/", pr.providerController.GetAllProviders)
	providerRouter.Get("/:id", pr.providerController.GetProvider)
	providerRouter.Post("/", pr.providerController.CreateProvider)
	providerRouter.Put("/:id", pr.providerController.UpdateProvider)
	providerRouter.Delete("/:id", pr.providerController.DeleteProvider)
}
