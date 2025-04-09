package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unxai/geonames-service/api"
	"github.com/unxai/geonames-service/db"
)

func main() {
	// 初始化数据库
	database := db.InitDB()
	defer database.Close()

	// 设置路由
	r := mux.NewRouter()

	// 注册路由
	api.RegisterRoutes(r)

	// 启动服务器
	log.Printf("服务器启动在 :8080 端口")
	log.Fatal(http.ListenAndServe(":8080", r))
}
