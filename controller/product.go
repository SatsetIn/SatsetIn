package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/watoken"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Handler to get all products
func GetAllProducts(respw http.ResponseWriter, r *http.Request) {
	data, err := atdb.GetAllDoc[[]model.Product](config.Mongoconn, "products", primitive.M{})
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Data kategori tidak ditemukan"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	if len(data) == 0 {
		var respn model.Response
		respn.Status = "Error: Data product kosong"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	var categories []map[string]interface{}
	for _, category := range data {
		categories = append(categories, map[string]interface{}{
			"id":             category.ID,
			"name":           category.Name,
			"description":    category.Description,
			"original_price": category.OriginalPrice,
			"discount_price": category.DiscountPrice,
			"image":          category.Image,
			"created_at":     category.CreatedAt,
			"updated_at":     category.UpdatedAt,
		})
	}

	at.WriteJSON(respw, http.StatusOK, categories)
}

// Handler to get a product by ID
func GetProductByID(w http.ResponseWriter, r *http.Request) {
	// Get the product ID from URL params
	id := r.URL.Query().Get("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Fetch the product from the database
	var product model.Product
	filter := bson.M{"_id": objectID}

	// Perbaiki pemanggilan GetOneDoc tanpa mengirimkan pointer
	product, err = atdb.GetOneDoc[model.Product](config.Mongoconn, "product", filter)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Respond with the product details
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// Handler to create a new product
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(r))
	if err != nil {
		_, err = watoken.Decode(config.PUBLICKEY, at.GetLoginFromHeader(r))
		if err != nil {
			var respn model.Response
			respn.Status = "Error: Token Tidak Valid"
			respn.Response = err.Error()
			at.WriteJSON(w, http.StatusForbidden, respn)
			return
		}
	}
	var product model.Product

	// Decode the request body into the Product struct
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Set the created and updated timestamps
	product.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	product.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	// Insert the product into the database
	_, err = atdb.InsertOneDoc(config.Mongoconn, "product", product)
	if err != nil {
		http.Error(w, "Error inserting product", http.StatusInternalServerError)
		return
	}

	// Respond with the created product
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "success", "message": "Product created"}
	json.NewEncoder(w).Encode(response)
}

// Handler to update an existing product
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Get the product ID from URL params
	id := r.URL.Query().Get("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Decode the request body into the Product struct
	var product model.Product
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Set the updated timestamp
	product.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	// Update the product in the database
	filter := bson.M{"_id": objectID}
	updateFields := bson.M{
		"name":           product.Name,
		"description":    product.Description,
		"discount_price": product.DiscountPrice,
		"original_price": product.OriginalPrice,
		"image":          product.Image,
		"updated_at":     product.UpdatedAt,
	}

	_, err = atdb.UpdateOneDoc(config.Mongoconn, "product", filter, updateFields)
	if err != nil {
		http.Error(w, "Error updating product", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "success", "message": "Product updated"}
	json.NewEncoder(w).Encode(response)
}

// Handler to delete a product
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Get the product ID from URL params
	id := r.URL.Query().Get("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Delete the product from the database
	filter := bson.M{"_id": objectID}
	_, err = atdb.DeleteOneDoc(config.Mongoconn, "product", filter)
	if err != nil {
		http.Error(w, "Error deleting product", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "success", "message": "Product deleted"}
	json.NewEncoder(w).Encode(response)
}
