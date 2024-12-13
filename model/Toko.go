package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Checkout struct {
	ID            primitive.ObjectID `bson:"_id,omitempty", json:"_id"`
	Address       string             `bson:"address,omitempty", json:"address"`
	Product       []Product          `bson:"product,omitempty", json:"product"`
	PaymentMethod string             `bson:"paymentmethod,omitempty", json:"paymentmethod"`
	TotalPrice    float64            `bson:"totalprice,omitempty", json:"totalprice"`
}
