package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unxai/geonames-service/db"
	"github.com/unxai/geonames-service/models"
)

// GetLocationsHandler 获取地理位置信息
func GetLocationsHandler(w http.ResponseWriter, r *http.Request) {
	db := db.GetDB()

	rows, err := db.Query("SELECT geoname_id, name, ascii_name, latitude, longitude, country_code, population, feature_class, feature_code FROM locations LIMIT 100")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var locations []models.Location
	for rows.Next() {
		var loc models.Location
		err := rows.Scan(
			&loc.GeonameID,
			&loc.Name,
			&loc.ASCII_Name,
			&loc.Latitude,
			&loc.Longitude,
			&loc.CountryCode,
			&loc.Population,
			&loc.FeatureClass,
			&loc.FeatureCode,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		locations = append(locations, loc)
	}

	json.NewEncoder(w).Encode(locations)
}

// GetLocationsByCountryHandler 按国家代码搜索
func GetLocationsByCountryHandler(w http.ResponseWriter, r *http.Request) {
	db := db.GetDB()

	vars := mux.Vars(r)
	countryCode := vars["countryCode"]

	rows, err := db.Query("SELECT geoname_id, name, ascii_name, latitude, longitude, country_code, population, feature_class, feature_code FROM locations WHERE country_code = $1", countryCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var locations []models.Location
	for rows.Next() {
		var loc models.Location
		err := rows.Scan(
			&loc.GeonameID,
			&loc.Name,
			&loc.ASCII_Name,
			&loc.Latitude,
			&loc.Longitude,
			&loc.CountryCode,
			&loc.Population,
			&loc.FeatureClass,
			&loc.FeatureCode,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		locations = append(locations, loc)
	}

	json.NewEncoder(w).Encode(locations)
}
