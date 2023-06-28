package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Password  string             `json:"-" bson:"password,omitempty"`
	Address   string             `json:"address,omitempty" bson:"address,omitempty"`
	Telephone string             `json:"telephone,omitempty" bson:"telephone,omitempty"`
	Role      Role               `json:"role,omitempty" bson:"role,omitempty"`
}

type Role struct {
	Name string `json:"name,omitempty" bson:"name,omitempty"`
}
