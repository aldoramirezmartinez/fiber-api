package main

import (
	"fmt"

	"github.com/aldoramirezmartinez/fiber-api/config"
	"github.com/aldoramirezmartinez/fiber-api/controllers"
	"github.com/aldoramirezmartinez/fiber-api/routes"
	"github.com/gofiber/fiber/v2"
)

type App struct {
	fiberApp                 *fiber.App
	UserController           *controllers.UserController
	ProviderController       *controllers.ProviderController
	PurchaseController       *controllers.PurchaseController
	ItemController           *controllers.ItemController
	PurchaseDetailController *controllers.PurchaseDetailController
}

func NewApp() *App {
	db, err := config.ConnectDB()
	if err != nil {
		fmt.Println("Failed to connect to MongoDB:", err)
	}

	userController := controllers.NewUserController(db)
	providerController := controllers.NewProviderController(db)
	purchaseController := controllers.NewPurchaseController(db)
	itemController := controllers.NewItemController(db)
	purchaseDetailController := controllers.NewPurchaseDetailController(db)

	fiberApp := fiber.New()

	return &App{
		fiberApp:                 fiberApp,
		UserController:           userController,
		ProviderController:       providerController,
		PurchaseController:       purchaseController,
		ItemController:           itemController,
		PurchaseDetailController: purchaseDetailController,
	}
}

func (app *App) Run() {
	userRoutes := routes.NewUserRoutes(app.fiberApp, app.UserController)
	userRoutes.SetupRoutes()

	providerRoutes := routes.NewProviderRoutes(app.fiberApp, app.ProviderController)
	providerRoutes.SetupRoutes()

	purchaseRoutes := routes.NewPurchaseRoutes(app.fiberApp, app.PurchaseController)
	purchaseRoutes.SetupRoutes()

	itemRoutes := routes.NewItemRoutes(app.fiberApp, app.ItemController)
	itemRoutes.SetupRoutes()

	purchaseDetailRoutes := routes.NewPurchaseDetailRoutes(app.fiberApp, app.PurchaseDetailController)
	purchaseDetailRoutes.SetupRoutes()

	port := config.GetPort()
	fmt.Println("Server listening on port:", port)

	err := app.fiberApp.Listen(":" + port)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
