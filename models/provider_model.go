package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Provider struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Address   string             `json:"address,omitempty" bson:"address,omitempty"`
	Telephone string             `json:"telephone,omitempty" bson:"telephone,omitempty"`
}
