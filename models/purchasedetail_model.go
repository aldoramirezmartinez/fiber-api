package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PurchaseDetail struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Quantity   int                `json:"quantity,omitempty" bson:"quantity,omitempty"`
	Total      float64            `json:"total,omitempty" bson:"total,omitempty"`
	ItemID     primitive.ObjectID `json:"item_id,omitempty" bson:"item_id,omitempty"`
	PurchaseID primitive.ObjectID `json:"purchase_id,omitempty" bson:"purchase_id,omitempty"`
}

type PurchaseDetailResponse struct {
	PurchaseDetail PurchaseDetail `json:"purchase_detail,omitempty" bson:"purchase_detail,omitempty"`
	Item           Item           `json:"item,omitempty" bson:"item,omitempty"`
	Purchase       Purchase       `json:"purchase,omitempty" bson:"purchase,omitempty"`
}
