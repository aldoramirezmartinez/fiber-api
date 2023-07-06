package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct {
	ID          primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Code        string             `json:"code,omitempty" bson:"code,omitempty"`
	UnitMeasure string             `json:"unit_measure,omitempty" bson:"unit_measure,omitempty"`
	Price       float64            `json:"price,omitempty" bson:"price,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	ProviderID  primitive.ObjectID `json:"-" bson:"provider_id,omitempty"`
}

type ItemResponse struct {
	Item     Item     `json:"item,omitempty" bson:"item,omitempty"`
	Provider Provider `json:"provider,omitempty" bson:"provider,omitempty"`
}
