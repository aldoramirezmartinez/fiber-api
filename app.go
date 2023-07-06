package main

import (
	"fmt"

	"github.com/aldoramirezmartinez/fiber-api/config"
	"github.com/aldoramirezmartinez/fiber-api/controllers"
	"github.com/aldoramirezmartinez/fiber-api/routes"
	"github.com/gofiber/fiber/v2"
)

type App struct {
	fiberApp             *fiber.App
	UserController       *controllers.UserController
	ProviderController   *controllers.ProviderController
	ItemController       *controllers.ItemController
	PurchaseV2Controller *controllers.PurchaseV2Controller
}

func NewApp() *App {
	db, err := config.ConnectDB()
	if err != nil {
		fmt.Println("Failed to connect to MongoDB:", err)
	}

	userController := controllers.NewUserController(db)
	providerController := controllers.NewProviderController(db)
	itemController := controllers.NewItemController(db)

	purchasev2Controller := controllers.NewPurchaseV2Controller(db)

	fiberApp := fiber.New()

	return &App{
		fiberApp:             fiberApp,
		UserController:       userController,
		ProviderController:   providerController,
		ItemController:       itemController,
		PurchaseV2Controller: purchasev2Controller,
	}
}

func (app *App) Run() {
	userRoutes := routes.NewUserRoutes(app.fiberApp, app.UserController)
	userRoutes.SetupRoutes()

	providerRoutes := routes.NewProviderRoutes(app.fiberApp, app.ProviderController)
	providerRoutes.SetupRoutes()

	itemRoutes := routes.NewItemRoutes(app.fiberApp, app.ItemController)
	itemRoutes.SetupRoutes()

	purchasev2Routes := routes.NewPurchaseV2Routes(app.fiberApp, app.PurchaseV2Controller)
	purchasev2Routes.SetupRoutes()

	port := config.GetPort()
	fmt.Println("Server listening on port:", port)

	err := app.fiberApp.Listen(":" + port)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
