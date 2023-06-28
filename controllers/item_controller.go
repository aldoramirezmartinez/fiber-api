package controllers

import (
	"context"

	"github.com/aldoramirezmartinez/fiber-api/models"
	"github.com/aldoramirezmartinez/fiber-api/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ItemController struct {
	itemCollection     *mongo.Collection
	providerCollection *mongo.Collection
}

func NewItemController(db *mongo.Database) *ItemController {
	return &ItemController{
		itemCollection:     db.Collection("items"),
		providerCollection: db.Collection("providers"),
	}
}

func (ic *ItemController) GetAllItems(c *fiber.Ctx) error {
	ctx := context.TODO()

	cursor, err := ic.itemCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get items",
			"error":   err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var items []models.Item
	if err := cursor.All(ctx, &items); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get items",
			"error":   err.Error(),
		})
	}

	var itemResponses []models.ItemResponse
	for _, item := range items {
		var provider models.Provider
		err := ic.providerCollection.FindOne(ctx, bson.M{"_id": item.ProviderID}).Decode(&provider)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to retrieve provider data",
				"error":   err.Error(),
			})
		}

		itemResponse := models.ItemResponse{
			Item:     item,
			Provider: provider,
		}
		itemResponses = append(itemResponses, itemResponse)
	}

	return c.JSON(itemResponses)
}

func (ic *ItemController) GetItem(c *fiber.Ctx) error {
	ctx := context.TODO()

	itemID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid item ID",
			"error":   err.Error(),
		})
	}

	var itemResponse models.ItemResponse

	err = ic.itemCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&itemResponse.Item)
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

	err = ic.providerCollection.FindOne(ctx, bson.M{"_id": itemResponse.Item.ProviderID}).Decode(&itemResponse.Provider)
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

	return c.JSON(itemResponse)
}

func (ic *ItemController) CreateItem(c *fiber.Ctx) error {
	ctx := context.TODO()

	item := new(models.Item)
	if err := c.BodyParser(item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to request body",
			"error":   err.Error(),
		})
	}

	providerID := item.ProviderID
	providerExists, err := utils.CheckDocumentExists(ic.providerCollection, providerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to check provider",
			"error":   err.Error(),
		})
	}
	if !providerExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Provider does not exist",
		})
	}

	result, err := ic.itemCollection.InsertOne(ctx, item)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create item",
			"error":   err.Error(),
		})
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get item ID",
			"error":   "Invalid inserted ID",
		})
	}

	item.ID = insertedID

	var provider models.Provider
	err = ic.providerCollection.FindOne(ctx, bson.M{"_id": item.ProviderID}).Decode(&provider)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve provider data",
			"error":   err.Error(),
		})
	}

	itemResponse := models.ItemResponse{
		Item:     *item,
		Provider: provider,
	}

	return c.JSON(itemResponse)
}

func (ic *ItemController) UpdateItem(c *fiber.Ctx) error {
	ctx := context.TODO()

	itemID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid item ID",
			"error":   err.Error(),
		})
	}

	itemToUpdate := new(models.Item)
	if err := c.BodyParser(itemToUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	var existingItem models.Item
	err = ic.itemCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&existingItem)
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

	providerID := itemToUpdate.ProviderID
	if providerID != existingItem.ProviderID {
		providerExists, err := utils.CheckDocumentExists(ic.providerCollection, providerID)
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

	itemToUpdate.ID = objID

	update := bson.M{
		"$set": itemToUpdate,
	}

	result, err := ic.itemCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update item",
			"error":   err.Error(),
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Item not found",
		})
	}

	var provider models.Provider
	err = ic.providerCollection.FindOne(ctx, bson.M{"_id": providerID}).Decode(&provider)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve provider data",
			"error":   err.Error(),
		})
	}

	itemResponse := models.ItemResponse{
		Item:     *itemToUpdate,
		Provider: provider,
	}

	return c.JSON(itemResponse)
}

func (ic *ItemController) DeleteItem(c *fiber.Ctx) error {
	ctx := context.TODO()

	itemID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid item ID",
			"error":   err.Error(),
		})
	}

	result, err := ic.itemCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete item",
			"error":   err.Error(),
		})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Item not found",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
