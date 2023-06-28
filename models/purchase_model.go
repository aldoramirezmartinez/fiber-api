package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Purchase struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	PurchaseOrder string             `json:"purchase_order,omitempty" bson:"purchase_order,omitempty"`
	Date          time.Time          `json:"date,omitempty" bson:"date,omitempty"`
	Status        string             `json:"status,omitempty" bson:"status,omitempty"`
	UserID        primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	ProviderID    primitive.ObjectID `json:"provider_id,omitempty" bson:"provider_id,omitempty"`
}

type PurchaseResponse struct {
	Purchase Purchase `json:"purchase,omitempty" bson:"purchase,omitempty"`
	User     User     `json:"user,omitempty" bson:"user,omitempty"`
	Provider Provider `json:"provider,omitempty" bson:"provider,omitempty"`
}
