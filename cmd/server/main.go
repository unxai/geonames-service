package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/unxai/geonames-service/api"
	"github.com/unxai/geonames-service/config"
	"github.com/unxai/geonames-service/db"
	"github.com/unxai/geonames-service/logger"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "geonames-server",
	Short: "GeoNames HTTP服务器",
	Long:  `GeoNames HTTP服务器，提供地理位置数据API服务。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runServer()
	},
}

func runServer() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 初始化日志
	if err := logger.InitLogger(cfg.Log.Level, cfg.Log.Path); err != nil {
		return fmt.Errorf("初始化日志失败: %w", err)
	}

	// 初始化数据库连接
	_ = db.GetDB()
	defer func() {
		if err := db.CloseDB(); err != nil {
			logger.Logger.Error("关闭数据库连接失败", zap.Error(err))
		}
	}()

	// 设置路由
	router := mux.NewRouter()
	api.RegisterRoutes(router)

	// 创建HTTP服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// 创建用于接收操作系统信号的通道
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动HTTP服务器
	go func() {
		logger.Logger.Info("服务器启动",
			zap.String("address", fmt.Sprintf("http://localhost%s", addr)),
		)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Logger.Error("服务器异常退出", zap.Error(err))
		}
	}()

	// 等待信号
	<-sigChan
	logger.Logger.Info("正在关闭服务器...")

	// 创建一个带超时的上下文用于优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("服务器关闭失败: %w", err)
	}

	logger.Logger.Info("服务器已关闭")
	return nil
}

func main() {
	// 加载配置文件
	_, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("加载配置文件失败: %v\n", err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		if logger.Logger != nil {
			logger.Logger.Error("执行命令失败", zap.Error(err))
		} else {
			fmt.Printf("执行命令失败: %v\n", err)
		}
		os.Exit(1)
	}
}
