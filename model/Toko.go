package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Toko struct {
	ID         primitive.ObjectID `bson:"_id,omitempty", json:"_id"`
	NamaToko   string             `bson:"nama_toko,omitempty", json:"nama_toko"`
	Categories string             `bson:"categories,omitempty", json:"categories"`
	Alamat     []Alamat           `json:"alamat"`
}

type Checkout struct {
	ID            primitive.ObjectID `bson:"_id,omitempty", json:"_id"`
	Address       string             `bson:"address,omitempty", json:"address"`
	Product       []Product          `bson:"product,omitempty", json:"product"`
	PaymentMethod string             `bson:"paymentmethod,omitempty", json:"paymentmethod"`
	TotalPrice    float64            `bson:"totalprice,omitempty", json:"totalprice"`
}
