package controllers

import (
	"context"
	"time"

	"github.com/aldoramirezmartinez/fiber-api/models"
	"github.com/aldoramirezmartinez/fiber-api/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PurchaseV2Controller struct {
	purchaseCollection *mongo.Collection
	userCollection     *mongo.Collection
	providerCollection *mongo.Collection
	itemCollection     *mongo.Collection
}

func NewPurchaseV2Controller(db *mongo.Database) *PurchaseV2Controller {
	return &PurchaseV2Controller{
		purchaseCollection: db.Collection("purchases"),
		userCollection:     db.Collection("users"),
		providerCollection: db.Collection("providers"),
		itemCollection:     db.Collection("items"),
	}
}

func (pc *PurchaseV2Controller) GetAllPurchasesV2(c *fiber.Ctx) error {
	ctx := context.TODO()

	// Obtener todas las compras de la versión 2 desde la base de datos
	cursor, err := pc.purchaseCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve purchases",
			"error":   err.Error(),
		})
	}
	defer cursor.Close(ctx)

	// Variable para almacenar las compras de la versión 2
	var purchases []models.PurchaseResponsev2

	// Iterar sobre el cursor y decodificar cada compra
	for cursor.Next(ctx) {
		var purchase models.Purchasev2
		if err := cursor.Decode(&purchase); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to decode purchase",
				"error":   err.Error(),
			})
		}

		// Obtener datos adicionales para la respuesta
		var user models.User
		err := pc.userCollection.FindOne(ctx, bson.M{"_id": purchase.UserID}).Decode(&user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to retrieve user data",
				"error":   err.Error(),
			})
		}

		var provider models.Provider
		err = pc.providerCollection.FindOne(ctx, bson.M{"_id": purchase.ProviderID}).Decode(&provider)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to retrieve provider data",
				"error":   err.Error(),
			})
		}

		// Crear respuesta de compra
		purchaseResponse := models.PurchaseResponsev2{
			Purchase: purchase,
			User:     user,
			Provider: provider,
		}

		for i := range purchaseResponse.Purchase.ItemList {
			itemID := purchaseResponse.Purchase.ItemList[i].ItemID

			var item models.Item
			err = pc.itemCollection.FindOne(ctx, bson.M{"_id": itemID}).Decode(&item)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Failed to retrieve item data",
					"error":   err.Error(),
				})
			}

			purchaseResponse.Purchase.ItemList[i].Item = item
		}

		// Agregar la compra a la lista de compras
		purchases = append(purchases, purchaseResponse)
	}

	if err := cursor.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve purchases",
			"error":   err.Error(),
		})
	}

	return c.JSON(purchases)
}

func (pc *PurchaseV2Controller) GetPurchaseV2(c *fiber.Ctx) error {
	ctx := context.TODO()

	purchaseOrder := c.Params("purchase_order")

	var purchase models.Purchasev2
	err := pc.purchaseCollection.FindOne(ctx, bson.M{"purchase_order": purchaseOrder}).Decode(&purchase)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Purchase not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve purchase",
			"error":   err.Error(),
		})
	}

	// Obtener datos adicionales para la respuesta
	var user models.User
	err = pc.userCollection.FindOne(ctx, bson.M{"_id": purchase.UserID}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve user data",
			"error":   err.Error(),
		})
	}

	var provider models.Provider
	err = pc.providerCollection.FindOne(ctx, bson.M{"_id": purchase.ProviderID}).Decode(&provider)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve provider data",
			"error":   err.Error(),
		})
	}

	// Crear respuesta de compra
	purchaseResponse := models.PurchaseResponsev2{
		Purchase: purchase,
		User:     user,
		Provider: provider,
	}

	for i := range purchaseResponse.Purchase.ItemList {
		itemID := purchaseResponse.Purchase.ItemList[i].ItemID

		var item models.Item
		err = pc.itemCollection.FindOne(ctx, bson.M{"_id": itemID}).Decode(&item)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to retrieve item data",
				"error":   err.Error(),
			})
		}

		purchaseResponse.Purchase.ItemList[i].Item = item
	}

	return c.JSON(purchaseResponse)
}

func (pc *PurchaseV2Controller) CreatePurchaseV2(c *fiber.Ctx) error {
	ctx := context.TODO()

	purchase := new(models.Purchasev2)
	if err := c.BodyParser(purchase); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse purchase body",
			"error":   err.Error(),
		})
	}

	userExists, err := pc.userCollection.CountDocuments(ctx, bson.M{"_id": purchase.UserID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to validate user",
			"error":   err.Error(),
		})
	}
	if userExists == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	providerExists, err := pc.providerCollection.CountDocuments(ctx, bson.M{"_id": purchase.ProviderID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to validate provider",
			"error":   err.Error(),
		})
	}
	if providerExists == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Provider not found",
		})
	}

	for i := range purchase.ItemList {
		itemExits, err := pc.itemCollection.CountDocuments(ctx, bson.M{"_id": purchase.ItemList[i].ItemID})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to validate item",
				"error":   err.Error(),
			})
		}
		if itemExits == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Item not found",
			})
		}

		var item models.Item
		err = pc.itemCollection.FindOne(ctx, bson.M{"_id": purchase.ItemList[i].ItemID}).Decode(&item)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"message": "Item not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to get item",
				"error":   err.Error(),
			})
		}

		subtotal := item.Price * float64(purchase.ItemList[i].Quantity)
		purchase.ItemList[i].Subtotal = subtotal
	}

	// Calcular el total
	var total float64
	for _, detail := range purchase.ItemList {
		total += detail.Subtotal
	}

	// Asignar valores al objeto de compra
	purchase.ID = primitive.NewObjectID()
	purchase.Date = time.Now()
	purchase.Total = total

	// Guardar la compra en la base de datos
	_, err = pc.purchaseCollection.InsertOne(ctx, purchase)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create purchase",
			"error":   err.Error(),
		})
	}

	// Obtener datos adicionales para la respuesta
	var user models.User
	err = pc.userCollection.FindOne(ctx, bson.M{"_id": purchase.UserID}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve user data",
			"error":   err.Error(),
		})
	}

	var provider models.Provider
	err = pc.providerCollection.FindOne(ctx, bson.M{"_id": purchase.ProviderID}).Decode(&provider)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve provider data",
			"error":   err.Error(),
		})
	}

	// Crear respuesta de compra
	purchaseResponse := models.PurchaseResponsev2{
		Purchase: *purchase,
		User:     user,
		Provider: provider,
	}

	for i := range purchaseResponse.Purchase.ItemList {
		itemID := purchaseResponse.Purchase.ItemList[i].ItemID

		var item models.Item
		err = pc.itemCollection.FindOne(ctx, bson.M{"_id": itemID}).Decode(&item)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to retrieve item data",
				"error":   err.Error(),
			})
		}

		purchaseResponse.Purchase.ItemList[i].Item = item

	}

	return c.JSON(purchaseResponse)
}

func (pvc *PurchaseV2Controller) UpdatePurchaseV2(c *fiber.Ctx) error {
	ctx := context.TODO()

	purchaseID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(purchaseID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid purchase ID",
			"error":   err.Error(),
		})
	}

	purchaseToUpdate := new(models.Purchasev2)
	if err := c.BodyParser(purchaseToUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	var existingPurchase models.Purchasev2
	err = pvc.purchaseCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&existingPurchase)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Purchase not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get purchase",
			"error":   err.Error(),
		})
	}

	purchaseToUpdate.Date = existingPurchase.Date

	userID := purchaseToUpdate.UserID
	if userID != existingPurchase.UserID {
		userExists, err := utils.CheckDocumentExists(pvc.userCollection, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to check user existence",
				"error":   err.Error(),
			})
		}
		if !userExists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "User does not exist",
			})
		}
	}

	providerID := purchaseToUpdate.ProviderID
	if providerID != existingPurchase.ProviderID {
		providerExists, err := utils.CheckDocumentExists(pvc.providerCollection, providerID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to check provider existence",
				"error":   err.Error(),
			})
		}
		if !providerExists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Provider does not exist",
			})
		}
	}

	_, err = pvc.purchaseCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": purchaseToUpdate})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update purchase",
			"error":   err.Error(),
		})
	}

	var user models.User
	err = pvc.userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve user data",
			"error":   err.Error(),
		})
	}

	var provider models.Provider
	err = pvc.providerCollection.FindOne(ctx, bson.M{"_id": providerID}).Decode(&provider)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve provider data",
			"error":   err.Error(),
		})
	}

	purchaseResponse := models.PurchaseResponsev2{
		Purchase: *purchaseToUpdate,
		User:     user,
		Provider: provider,
	}

	return c.JSON(purchaseResponse)
}
