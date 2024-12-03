package model

import (
	"context"
	"net/http"
	"time"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/watoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Definisikan tipe Response untuk hasil API
type Response struct {
	Response string `json:"response"`
	Info     string `json:"info,omitempty"`
	Status   string `json:"status,omitempty"`
	Location string `json:"location,omitempty"`
}

// Definisikan struktur User
type User struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty"`
	Email      string             `json:"email,omitempty" bson:"email,omitempty"`
	Password   string             `json:"password,omitempty" bson:"password,omitempty"`
	Phonenumber string            `json:"phonenumber,omitempty" bson:"phonenumber,omitempty"` // pastikan ada field phonenumber
}

// Fungsi untuk membuat pengguna baru
func CreateUser(user *User, collection *mongo.Collection) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return collection.InsertOne(ctx, user)
}

// Fungsi untuk mencari pengguna berdasarkan email
func FindUserByEmail(email string, collection *mongo.Collection) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return collection.FindOne(ctx, bson.M{"email": email})
}

// Fungsi untuk mendapatkan data pengguna berdasarkan token yang diterima dalam request
func GetDataUser(respw http.ResponseWriter, req *http.Request) {
	// Decode token dari header
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		var respn Response
		respn.Status = "Error: Token Tidak Valid"
		respn.Info = config.PublicKeyWhatsAuth
		respn.Location = "Decode Token Error: " + at.GetLoginFromHeader(req)
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusForbidden, respn)
		return
	}

	// Ambil data pengguna dari database MongoDB berdasarkan phonenumber
	docuser, err := atdb.GetOneDoc[User](config.Mongoconn, "user", primitive.M{"phonenumber": payload.Id})
	if err != nil {
		// Jika pengguna tidak ditemukan, kirim respons dengan informasi pengguna dari token
		var newUser User
		newUser.Phonenumber = payload.Id
		newUser.Name = payload.Alias
		at.WriteJSON(respw, http.StatusNotFound, newUser)
		return
	}

	// Update data pengguna dengan informasi dari token
	docuser.Name = payload.Alias
	at.WriteJSON(respw, http.StatusOK, docuser)
}
