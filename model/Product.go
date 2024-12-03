package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Product represents a product in the database
type Product struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name          string             `bson:"name" json:"name"`
	Description   string             `bson:"description" json:"description"`
	DiscountPrice int64              `bson:"discount_price" json:"discount_price"`
	OriginalPrice int64              `bson:"original_price" json:"original_price"`
	Image         string             `bson:"image" json:"image"`
	CreatedAt     primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt     primitive.DateTime `json:"updated_at" bson:"updated_at"`
}
