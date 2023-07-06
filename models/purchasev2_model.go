package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Purchasev2 struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	PurchaseOrder string             `json:"purchase_order,omitempty" bson:"purchase_order,omitempty"`
	Date          time.Time          `json:"date,omitempty" bson:"date,omitempty"`
	Status        string             `json:"status,omitempty" bson:"status,omitempty"`
	ItemList      []PurchaseDetailv2 `json:"item_list,omitempty" bson:"item_list,omitempty"`
	Total         float64            `json:"total,omitempty" bson:"total,omitempty"`
	UserID        primitive.ObjectID `json:"-" bson:"user_id,omitempty"`
	ProviderID    primitive.ObjectID `json:"-" bson:"provider_id,omitempty"`
}

type PurchaseDetailv2 struct {
	ItemID   primitive.ObjectID `json:"item_id,omitempty" bson:"item_id,omitempty"`
	Item     Item               `json:"item,omitempty" bson:"item,omitempty"`
	Quantity int                `json:"quantity,omitempty" bson:"quantity,omitempty"`
	Subtotal float64            `json:"subtotal,omitempty" bson:"subtotal,omitempty"`
}

type PurchaseResponsev2 struct {
	Purchase Purchasev2 `json:"purchase,omitempty" bson:"purchase,omitempty"`
	User     User       `json:"user,omitempty" bson:"user,omitempty"`
	Provider Provider   `json:"provider,omitempty" bson:"provider,omitempty"`
}
