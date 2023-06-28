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

type PurchaseController struct {
	purchaseCollection *mongo.Collection
	userCollection     *mongo.Collection
	providerCollection *mongo.Collection
}

func NewPurchaseController(db *mongo.Database) *PurchaseController {
	return &PurchaseController{
		purchaseCollection: db.Collection("purchases"),
		userCollection:     db.Collection("users"),
		providerCollection: db.Collection("providers"),
	}
}

func (pc *PurchaseController) GetAllPurchases(c *fiber.Ctx) error {
	ctx := context.TODO()

	cursor, err := pc.purchaseCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get purchases",
			"error":   err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var purchases []models.Purchase
	if err := cursor.All(ctx, &purchases); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get purchases",
			"error":   err.Error(),
		})
	}

	var purchaseResponses []models.PurchaseResponse
	for _, purchase := range purchases {
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

		purchaseResponse := models.PurchaseResponse{
			Purchase: purchase,
			User:     user,
			Provider: provider,
		}
		purchaseResponses = append(purchaseResponses, purchaseResponse)
	}

	return c.JSON(purchaseResponses)
}

func (pc *PurchaseController) GetPurchase(c *fiber.Ctx) error {
	ctx := context.TODO()

	purchaseID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(purchaseID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid purchase ID",
			"error":   err.Error(),
		})
	}

	var purchaseResponse models.PurchaseResponse

	err = pc.purchaseCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&purchaseResponse.Purchase)
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

	err = pc.userCollection.FindOne(ctx, bson.M{"_id": purchaseResponse.Purchase.UserID}).Decode(&purchaseResponse.User)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get user",
			"error":   err.Error(),
		})
	}

	err = pc.providerCollection.FindOne(ctx, bson.M{"_id": purchaseResponse.Purchase.ProviderID}).Decode(&purchaseResponse.Provider)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Provider not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get provider",
			"error":   err.Error(),
		})
	}

	return c.JSON(purchaseResponse)
}

func (pc *PurchaseController) CreatePurchase(c *fiber.Ctx) error {
	ctx := context.TODO()

	purchase := new(models.Purchase)
	if err := c.BodyParser(purchase); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	purchase.Date = time.Now()

	userID := purchase.UserID
	userExists, err := utils.CheckDocumentExists(pc.userCollection, userID)
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

	providerID := purchase.ProviderID
	providerExists, err := utils.CheckDocumentExists(pc.providerCollection, providerID)
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

	result, err := pc.purchaseCollection.InsertOne(ctx, purchase)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create purchase",
			"error":   err.Error(),
		})
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create purchase",
			"error":   "Invalid inserted ID",
		})
	}

	purchase.ID = insertedID

	var user models.User
	err = pc.userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve user data",
			"error":   err.Error(),
		})
	}

	var provider models.Provider
	err = pc.providerCollection.FindOne(ctx, bson.M{"_id": providerID}).Decode(&provider)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve provider data",
			"error":   err.Error(),
		})
	}

	purchaseResponse := models.PurchaseResponse{
		Purchase: *purchase,
		User:     user,
		Provider: provider,
	}

	return c.JSON(purchaseResponse)
}

func (pc *PurchaseController) UpdatePurchase(c *fiber.Ctx) error {
	ctx := context.TODO()

	purchaseID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(purchaseID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid purchase ID",
			"error":   err.Error(),
		})
	}

	purchaseToUpdate := new(models.Purchase)
	if err := c.BodyParser(purchaseToUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	var existingPurchase models.Purchase
	err = pc.purchaseCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&existingPurchase)
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
		userExists, err := utils.CheckDocumentExists(pc.userCollection, userID)
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
		providerExists, err := utils.CheckDocumentExists(pc.providerCollection, providerID)
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

	purchaseToUpdate.ID = objID

	update := bson.M{
		"$set": purchaseToUpdate,
	}

	result, err := pc.purchaseCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update purchase",
			"error":   err.Error(),
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Purchase not found",
		})
	}

	var user models.User
	err = pc.userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve user data",
			"error":   err.Error(),
		})
	}

	var provider models.Provider
	err = pc.providerCollection.FindOne(ctx, bson.M{"_id": providerID}).Decode(&provider)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve provider data",
			"error":   err.Error(),
		})
	}

	purchaseResponse := models.PurchaseResponse{
		Purchase: *purchaseToUpdate,
		User:     user,
		Provider: provider,
	}

	return c.JSON(purchaseResponse)
}

func (pc *PurchaseController) DeletePurchase(c *fiber.Ctx) error {
	ctx := context.TODO()

	purchaseID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(purchaseID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid purchase ID",
			"error":   err.Error(),
		})
	}

	result, err := pc.purchaseCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete purchase",
			"error":   err.Error(),
		})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Purchase not found",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
