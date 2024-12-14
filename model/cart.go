package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CartItem represents an item in the cart
type CartItem struct {
	ProductID primitive.ObjectID `json:"product_id,omitempty" bson:"product_id,omitempty"`
	Quantity  int                `json:"quantity,omitempty" bson:"quantity,omitempty"`
}

// Cart represents a shopping cart
type Cart struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	User      []User             `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Items     []CartItem         `json:"items,omitempty" bson:"items,omitempty"`
	NamaToko  string             `json:"nama_toko,omitempty" bson:"nama_toko,omitempty"`
	Category  string             `json:"category,omitempty" bson:"category,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	Alamat    []Alamat           `json:"alamat"`
}

type Alamat struct {
	Street      string `json:"street,omitempty" bson:"street,omitempty"`
	Province    string `json:"province,omitempty" bson:"province,omitempty"`
	City        string `json:"city,omitempty" bson:"city,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	PostalCode  string `json:"postal_code,omitempty" bson:"postal_code,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

// AddItem adds an item to the cart
func (c *Cart) AddItem(item CartItem) {
	// Check if item already exists, then update the quantity
	for i, cartItem := range c.Items {
		if cartItem.ProductID == item.ProductID {
			c.Items[i].Quantity += item.Quantity
			return
		}
	}
	// If item doesn't exist, append it to the cart
	c.Items = append(c.Items, item)
}

// RemoveItem removes an item from the cart by its ProductID
func (c *Cart) RemoveItem(productID primitive.ObjectID) {
	for i, cartItem := range c.Items {
		if cartItem.ProductID == productID {
			// Remove the item from the slice
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			return
		}
	}
}

// TotalPrice calculates the total price of the cart by summing up prices of all items
func (c *Cart) TotalPrice(productsCollection *mongo.Collection) (int64, error) {
	var totalPrice int64

	// Loop through each cart item and calculate the total price
	for _, cartItem := range c.Items {
		// Fetch the product from the products collection using ProductID
		var product Product
		err := productsCollection.FindOne(context.Background(), bson.M{"_id": cartItem.ProductID}).Decode(&product)
		if err != nil {
			return 0, err // Return error if product is not found
		}

		// Determine which price to use (DiscountPrice if available, else OriginalPrice)
		var price int64
		if product.DiscountPrice > 0 {
			price = product.DiscountPrice
		} else {
			price = product.OriginalPrice
		}

		// Add the total price of the item (price * quantity)
		totalPrice += price * int64(cartItem.Quantity)
	}

	// Return the calculated total price
	return totalPrice, nil
}
