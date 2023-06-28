package controllers

import (
	"context"

	"github.com/aldoramirezmartinez/fiber-api/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProviderController struct {
	collection *mongo.Collection
}

func NewProviderController(db *mongo.Database) *ProviderController {
	return &ProviderController{
		collection: db.Collection("providers"),
	}
}

func (pc *ProviderController) GetAllProviders(c *fiber.Ctx) error {
	ctx := context.TODO()

	cursor, err := pc.collection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get providers",
			"error":   err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var providers []models.Provider
	if err := cursor.All(ctx, &providers); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get providers",
			"error":   err.Error(),
		})
	}

	return c.JSON(providers)
}

func (pc *ProviderController) GetProvider(c *fiber.Ctx) error {
	ctx := context.TODO()

	providerID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(providerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid provider ID",
			"error":   err.Error(),
		})
	}

	var provider models.Provider
	err = pc.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&provider)
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

	return c.JSON(provider)
}

func (pc *ProviderController) CreateProvider(c *fiber.Ctx) error {
	ctx := context.TODO()

	provider := new(models.Provider)
	if err := c.BodyParser(provider); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	result, err := pc.collection.InsertOne(ctx, provider)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create provider",
			"error":   err.Error(),
		})
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create provider",
		})
	}

	provider.ID = insertedID

	return c.JSON(provider)
}

func (pc *ProviderController) UpdateProvider(c *fiber.Ctx) error {
	ctx := context.TODO()

	providerID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(providerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid provider ID",
			"error":   err.Error(),
		})
	}

	updateData := new(models.Provider)
	if err := c.BodyParser(updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	updateData.ID = objID

	update := bson.M{
		"$set": updateData,
	}

	result, err := pc.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update provider",
			"error":   err.Error(),
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Provider not found",
		})
	}

	return c.JSON(updateData)
}

func (pc *ProviderController) DeleteProvider(c *fiber.Ctx) error {
	ctx := context.TODO()

	providerID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(providerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid provider ID",
			"error":   err.Error(),
		})
	}

	result, err := pc.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete provider",
			"error":   err.Error(),
		})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Provider not found",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
