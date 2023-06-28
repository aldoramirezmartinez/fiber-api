package controllers

import (
	"context"

	"github.com/aldoramirezmartinez/fiber-api/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	collection *mongo.Collection
}

func NewUserController(db *mongo.Database) *UserController {
	return &UserController{
		collection: db.Collection("users"),
	}
}

func (uc *UserController) GetAllUsers(c *fiber.Ctx) error {
	ctx := context.TODO()

	cursor, err := uc.collection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get users",
			"error":   err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get users",
			"error":   err.Error(),
		})
	}

	return c.JSON(users)
}

func (uc *UserController) GetUser(c *fiber.Ctx) error {
	ctx := context.TODO()

	userID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid user ID",
			"error":   err.Error(),
		})
	}

	var user models.User
	err = uc.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
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

	return c.JSON(user)
}

func (uc *UserController) CreateUser(c *fiber.Ctx) error {
	ctx := context.TODO()

	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	result, err := uc.collection.InsertOne(ctx, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user",
			"error":   err.Error(),
		})
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user",
		})
	}

	user.ID = insertedID

	return c.JSON(user)
}

func (uc *UserController) UpdateUser(c *fiber.Ctx) error {
	ctx := context.TODO()

	userID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid user ID",
			"error":   err.Error(),
		})
	}

	updateData := new(models.User)
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

	result, err := uc.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update user",
			"error":   err.Error(),
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	return c.JSON(updateData)
}

func (uc *UserController) DeleteUser(c *fiber.Ctx) error {
	ctx := context.TODO()

	userID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid user ID",
			"error":   err.Error(),
		})
	}

	result, err := uc.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete user",
			"error":   err.Error(),
		})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
