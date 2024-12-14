package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"

	// "github.com/gocroot/helper/watoken"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
)

func GetRegion(respw http.ResponseWriter, req *http.Request) {

	// Parse koordinat dari body request
	var longlat model.LongLat
	json.NewDecoder(req.Body).Decode(&longlat)
	// Filter query geospasial
	filter := bson.M{
		"border": bson.M{
			"$geoIntersects": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{longlat.Longitude, longlat.Latitude},
				},
			},
		},
	}

	// Cari region berdasarkan filter
	region, err := atdb.GetOneDoc[model.Region](config.MongoconnGeo, "region", filter)
	if err != nil {
		at.WriteJSON(respw, http.StatusNotFound, bson.M{"error": "Region not found"})
		return
	}

	// Format respon sebagai FeatureCollection GeoJSON
	geoJSON := bson.M{
		"type": "FeatureCollection",
		"features": []bson.M{
			{
				"type": "Feature",
				"geometry": bson.M{
					"type":        region.Border.Type,
					"coordinates": region.Border.Coordinates,
				},
				"properties": bson.M{
					"province":     region.Province,
					"district":     region.District,
					"sub_district": region.SubDistrict,
					"village":      region.Village,
				},
			},
		},
	}

	// Kirim respon dalam format GeoJSON
	at.WriteJSON(respw, http.StatusOK, geoJSON)
}

// GET ROADS
func GetRoads(respw http.ResponseWriter, req *http.Request) {
	var longlat model.LongLat
	if err := json.NewDecoder(req.Body).Decode(&longlat); err != nil {
		http.Error(respw, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validasi data input
	if longlat.Longitude == 0 || longlat.Latitude == 0 || longlat.MaxDistance <= 0 {
		http.Error(respw, "Invalid longitude, latitude, or max distance", http.StatusBadRequest)
		return
	}

	// Query MongoDB untuk filter geospasial
	filter := bson.M{
		"geometry": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{longlat.Longitude, longlat.Latitude},
				},
				"$maxDistance": longlat.MaxDistance,
			},
		},
	}

	// Ambil dokumen dari koleksi MongoDB
	roads, err := atdb.GetAllDoc[[]model.Roads](config.MongoconnGeo, "roads", filter)
	if err != nil {
		log.Printf("Error fetching roads: %v", err)
		at.WriteJSON(respw, http.StatusNotFound, nil)
		return
	}

	// Kirimkan data sebagai respon
	at.WriteJSON(respw, http.StatusOK, roads)
}

