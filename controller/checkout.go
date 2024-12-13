package controller

import (
	// "context"
	// "fmt"
	"encoding/json"
	"net/http"

	// "time"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/watoken"
	"github.com/gocroot/model"
)

func Createcheckout(respw http.ResponseWriter, req *http.Request) {
	// Decode token
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		_, err = watoken.Decode(config.PUBLICKEY, at.GetLoginFromHeader(req))
		if err != nil {
			var respn model.Response
			respn.Status = "Error: Token Tidak Valid"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusForbidden, respn)
			return
		}
	}

	var checkout model.Checkout
	if err := json.NewDecoder(req.Body).Decode(&checkout); err != nil {
		respn := model.Response{
			Status:   "Invalid Request",
			Response: err.Error(),
		}
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	newCheckout := model.Checkout{
		Address:       checkout.Address,
		Product:       checkout.Product,
		PaymentMethod: checkout.PaymentMethod,
		TotalPrice:    checkout.TotalPrice,
	}

	_, err = atdb.InsertOneDoc(config.Mongoconn, "checkout", newCheckout)
	if err != nil {
		respn := model.Response{
			Status:   "Failed to insert new user",
			Response: err.Error(),
		}
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	response := map[string]interface{}{
		"status":       "success",
		"message":      "checkout berhasil di buat",
		"nama_product": newCheckout.Product,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}
