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

type PurchaseDetailController struct {
	purchaseDetailCollection *mongo.Collection
	itemCollection           *mongo.Collection
	purchaseCollection       *mongo.Collection
}

func NewPurchaseDetailController(db *mongo.Database) *PurchaseDetailController {
	return &PurchaseDetailController{
		purchaseDetailCollection: db.Collection("purchase_details"),
		itemCollection:           db.Collection("items"),
		purchaseCollection:       db.Collection("purchases"),
	}
}

func (pdc *PurchaseDetailController) GetAllPurchaseDetails(c *fiber.Ctx) error {
	ctx := context.TODO()

	cursor, err := pdc.purchaseDetailCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get purchase details",
			"error":   err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var purchaseDetails []models.PurchaseDetail
	if err := cursor.All(ctx, &purchaseDetails); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get purchase details",
			"error":   err.Error(),
		})
	}

	var purchaseDetailResponses []models.PurchaseDetailResponse
	for _, purchaseDetail := range purchaseDetails {
		var item models.Item
		err := pdc.itemCollection.FindOne(ctx, bson.M{"_id": purchaseDetail.ItemID}).Decode(&item)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to retrieve item data",
				"error":   err.Error(),
			})
		}

		var purchase models.Purchase
		err = pdc.purchaseCollection.FindOne(ctx, bson.M{"_id": purchaseDetail.PurchaseID}).Decode(&purchase)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to retrieve purchase data",
				"error":   err.Error(),
			})
		}

		purchaseDetailResponse := models.PurchaseDetailResponse{
			PurchaseDetail: purchaseDetail,
			Item:           item,
			Purchase:       purchase,
		}
		purchaseDetailResponses = append(purchaseDetailResponses, purchaseDetailResponse)
	}

	return c.JSON(purchaseDetailResponses)
}

func (pdc *PurchaseDetailController) GetPurchaseDetail(c *fiber.Ctx) error {
	ctx := context.TODO()

	purchaseDetailID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(purchaseDetailID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid purchase detail ID",
			"error":   err.Error(),
		})
	}

	var purchaseDetailResponse models.PurchaseDetailResponse

	err = pdc.purchaseDetailCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&purchaseDetailResponse.PurchaseDetail)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Purchase detail not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get purchase detail",
			"error":   err.Error(),
		})
	}

	err = pdc.itemCollection.FindOne(ctx, bson.M{"_id": purchaseDetailResponse.PurchaseDetail.ItemID}).Decode(&purchaseDetailResponse.Item)
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

	err = pdc.purchaseCollection.FindOne(ctx, bson.M{"_id": purchaseDetailResponse.PurchaseDetail.PurchaseID}).Decode(&purchaseDetailResponse.Purchase)
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

	return c.JSON(purchaseDetailResponse)
}

func (pdc *PurchaseDetailController) CreatePurchaseDetail(c *fiber.Ctx) error {
	ctx := context.TODO()

	purchaseDetail := new(models.PurchaseDetail)
	if err := c.BodyParser((purchaseDetail)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	itemID := purchaseDetail.ItemID
	itemExists, err := utils.CheckDocumentExists(pdc.itemCollection, itemID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to check item existence",
			"error":   err.Error(),
		})
	}
	if !itemExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Item does not exist",
		})
	}

	purchaseID := purchaseDetail.PurchaseID
	purchaseExists, err := utils.CheckDocumentExists(pdc.purchaseCollection, purchaseID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to check purchase existence",
			"error":   err.Error(),
		})
	}
	if !purchaseExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Purchase does not exist",
		})
	}

	var item models.Item
	err = pdc.itemCollection.FindOne(ctx, bson.M{"_id": itemID}).Decode(&item)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve item data",
			"error":   err.Error(),
		})
	}

	var purchase models.Purchase
	err = pdc.purchaseCollection.FindOne(ctx, bson.M{"_id": purchaseID}).Decode(&purchase)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve purchase data",
			"error":   err.Error(),
		})
	}

	if item.ProviderID != purchase.ProviderID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Item and Purchase have different providers",
		})
	}

	purchaseDetail.Total = float64(purchaseDetail.Quantity) * item.Price

	result, err := pdc.purchaseDetailCollection.InsertOne(ctx, purchaseDetail)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create purchase detail",
			"error":   err.Error(),
		})
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create purchase detail",
			"error":   "Invalid inserted ID",
		})
	}

	purchaseDetail.ID = insertedID

	purchaseDetailResponse := models.PurchaseDetailResponse{
		PurchaseDetail: *purchaseDetail,
		Item:           item,
		Purchase:       purchase,
	}

	return c.JSON(purchaseDetailResponse)
}

func (pdc *PurchaseDetailController) UpdatePurchaseDetail(c *fiber.Ctx) error {
	ctx := context.TODO()

	purchaseDetailID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(purchaseDetailID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid purchase detail ID",
			"error":   err.Error(),
		})
	}

	purchaseDetailToUpdate := new(models.PurchaseDetail)
	if err := c.BodyParser(purchaseDetailToUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	var existingPurchaseDetail models.PurchaseDetail
	err = pdc.purchaseDetailCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&existingPurchaseDetail)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Purchase detail not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get purchase detail",
			"error":   err.Error(),
		})
	}

	existingPurchaseDetail.Quantity = purchaseDetailToUpdate.Quantity

	var item models.Item
	err = pdc.itemCollection.FindOne(ctx, bson.M{"_id": existingPurchaseDetail.ItemID}).Decode(&item)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve item data",
			"error":   err.Error(),
		})
	}

	existingPurchaseDetail.Total = float64(existingPurchaseDetail.Quantity) * item.Price

	update := bson.M{
		"$set": bson.M{
			"quantity": existingPurchaseDetail.Quantity,
			"total":    existingPurchaseDetail.Total,
		},
	}

	result, err := pdc.purchaseDetailCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update purchase detail",
			"error":   err.Error(),
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Purchase detail not found",
		})
	}

	var purchase models.Purchase
	err = pdc.purchaseCollection.FindOne(ctx, bson.M{"_id": existingPurchaseDetail.PurchaseID}).Decode(&purchase)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve purchase data",
			"error":   err.Error(),
		})
	}

	purchaseDetailResponse := models.PurchaseDetailResponse{
		PurchaseDetail: existingPurchaseDetail,
		Item:           item,
		Purchase:       purchase,
	}

	return c.JSON(purchaseDetailResponse)
}

func (pdc *PurchaseDetailController) DeletePurchaseDetail(c *fiber.Ctx) error {
	ctx := context.TODO()

	purchaseDetailID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(purchaseDetailID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid purchase detail ID",
			"error":   err.Error(),
		})
	}

	result, err := pdc.purchaseDetailCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete purchase detail",
			"error":   err.Error(),
		})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Purchase detail not found",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
