package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserParse struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	UserID   string             `bson:"user_id"`
	Name     string             `bson:"name"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

type ProductParse struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ProductID string             `json:"product_id"`
	Name      string             `json:"name"`
	Price     int32              `json:"price"`
	Stock     int32              `json:"stock"`
}
