package api

import (
	"github.com/gorilla/mux"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *mux.Router) {
	// 获取地理位置信息
	r.HandleFunc("/locations", GetLocationsHandler).Methods("GET")

	// 按国家代码搜索
	r.HandleFunc("/locations/{countryCode}", GetLocationsByCountryHandler).Methods("GET")
}
