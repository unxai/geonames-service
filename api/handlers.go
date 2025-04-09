package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unxai/geonames-service/models"
	"github.com/unxai/geonames-service/utils"
)

var locations []models.Location

// DownloadHandler 处理下载地理数据的请求
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	locations, err = utils.DownloadGeoData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "数据下载成功",
		"count":   len(locations),
	}
	json.NewEncoder(w).Encode(response)
}

// GetLocationsHandler 获取地理位置信息
func GetLocationsHandler(w http.ResponseWriter, r *http.Request) {
	if len(locations) == 0 {
		http.Error(w, "请先下载数据", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(locations[:100]) // 仅返回前100条记录作为示例
}

// GetLocationsByCountryHandler 按国家代码搜索
func GetLocationsByCountryHandler(w http.ResponseWriter, r *http.Request) {
	if len(locations) == 0 {
		http.Error(w, "请先下载数据", http.StatusNotFound)
		return
	}

	vars := mux.Vars(r)
	countryCode := vars["countryCode"]

	var result []models.Location
	for _, loc := range locations {
		if loc.CountryCode == countryCode {
			result = append(result, loc)
		}
	}

	json.NewEncoder(w).Encode(result)
}
