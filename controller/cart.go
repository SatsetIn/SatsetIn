package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateCart - Fungsi untuk membuat cart baru
func CreateCart(w http.ResponseWriter, r *http.Request) {
	var cart model.Cart

	// Decode JSON request ke struct Cart
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&cart); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Generate timestamp untuk CreatedAt dan UpdatedAt
	cart.CreatedAt = time.Now()
	cart.UpdatedAt = time.Now()

	// Insert Cart ke dalam database
	insertedID, err := atdb.InsertOneDoc(config.Mongoconn, "cart", cart)
	if err != nil {
		http.Error(w, "Failed to create cart", http.StatusInternalServerError)
		return
	}

	// Kirim response kembali ke client
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":  "success",
		"cart_id": insertedID,
	}
	json.NewEncoder(w).Encode(response)
}

// GetCart - Fungsi untuk mendapatkan cart berdasarkan user_id
func GetCart(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan user_id dari URL params
	userID := r.URL.Query().Get("user_id")
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Mencari cart berdasarkan user_id
	filter := bson.M{"user_id": objectID}

	// Mendapatkan cart menggunakan GetOneDoc
	cart, err := atdb.GetOneDoc[model.Cart](config.Mongoconn, "cart", filter)
	if err != nil {
		http.Error(w, "Cart not found", http.StatusNotFound)
		return
	}

	// Mengirimkan response cart ke client dengan format JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cart); err != nil {
		http.Error(w, "Error encoding cart", http.StatusInternalServerError)
		return
	}
}

func AddItemToCart(w http.ResponseWriter, r *http.Request) {
	var cartItem model.CartItem

	// Decode JSON request ke struct CartItem
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&cartItem); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Mendapatkan user_id dari URL params
	userID := r.URL.Query().Get("user_id")
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Mencari cart berdasarkan user_id
	filter := bson.M{"user_id": objectID}
	cart, err := atdb.GetOneDoc[model.Cart](config.Mongoconn, "cart", filter)
	if err != nil {
		http.Error(w, "Cart not found", http.StatusNotFound)
		return
	}

	// Menambahkan item ke dalam cart
	cart.Items = append(cart.Items, cartItem)
	cart.UpdatedAt = time.Now()

	// Update cart di database
	updateFields := bson.M{"items": cart.Items, "updated_at": cart.UpdatedAt}
	_, err = atdb.UpdateOneDoc(config.Mongoconn, "cart", filter, updateFields)
	if err != nil {
		http.Error(w, "Failed to add item to cart", http.StatusInternalServerError)
		return
	}

	// Mengirimkan response ke client
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "success", "message": "Item added to cart"}
	json.NewEncoder(w).Encode(response)
}

// UpdateItemInCart - Fungsi untuk memperbarui kuantitas item dalam cart
func UpdateItemInCart(w http.ResponseWriter, r *http.Request) {
	var cartItem model.CartItem

	// Decode JSON request ke struct CartItem
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&cartItem); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Mendapatkan user_id dan product_id dari URL params
	userID := r.URL.Query().Get("user_id")
	productID := r.URL.Query().Get("product_id")
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	productObjectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Mencari cart berdasarkan user_id
	filter := bson.M{"user_id": objectID}
	cart, err := atdb.GetOneDoc[model.Cart](config.Mongoconn, "cart", filter)
	if err != nil {
		http.Error(w, "Cart not found", http.StatusNotFound)
		return
	}

	// Mencari item berdasarkan product_id dan memperbarui kuantitasnya
	itemFound := false
	for i, item := range cart.Items {
		if item.ProductID == productObjectID {
			cart.Items[i].Quantity = cartItem.Quantity
			itemFound = true
			break
		}
	}

	if !itemFound {
		http.Error(w, "Product not found in cart", http.StatusNotFound)
		return
	}

	// Update waktu terakhir diubah
	cart.UpdatedAt = time.Now()

	// Update cart di database
	updateFields := bson.M{"items": cart.Items, "updated_at": cart.UpdatedAt}
	_, err = atdb.UpdateOneDoc(config.Mongoconn, "cart", filter, updateFields)
	if err != nil {
		http.Error(w, "Failed to update item in cart", http.StatusInternalServerError)
		return
	}

	// Mengirimkan response ke client
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "success", "message": "Item updated in cart"}
	json.NewEncoder(w).Encode(response)
}

// DeleteItemFromCart - Fungsi untuk menghapus item dari cart
func DeleteItemFromCart(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan user_id dan product_id dari URL params
	userID := r.URL.Query().Get("user_id")
	productID := r.URL.Query().Get("product_id")
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	productObjectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Mencari cart berdasarkan user_id
	filter := bson.M{"user_id": objectID}
	cart, err := atdb.GetOneDoc[model.Cart](config.Mongoconn, "cart", filter)
	if err != nil {
		http.Error(w, "Cart not found", http.StatusNotFound)
		return
	}

	// Menghapus item dari cart
	var updatedItems []model.CartItem
	for _, item := range cart.Items {
		if item.ProductID != productObjectID {
			updatedItems = append(updatedItems, item)
		}
	}

	cart.Items = updatedItems
	cart.UpdatedAt = time.Now()

	// Update cart di database
	updateFields := bson.M{"items": cart.Items, "updated_at": cart.UpdatedAt}
	_, err = atdb.UpdateOneDoc(config.Mongoconn, "cart", filter, updateFields)
	if err != nil {
		http.Error(w, "Failed to delete item from cart", http.StatusInternalServerError)
		return
	}

	// Mengirimkan response ke client
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "success", "message": "Item deleted from cart"}
	json.NewEncoder(w).Encode(response)
}
