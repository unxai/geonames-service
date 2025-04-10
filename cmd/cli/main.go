package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/unxai/geonames-service/config"
	"github.com/unxai/geonames-service/db"
	"github.com/unxai/geonames-service/logger"
	"github.com/unxai/geonames-service/utils"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "geonames-cli",
	Short: "GeoNames CLI工具",
	Long:  `GeoNames CLI工具用于管理GeoNames数据。`,
}

// 下载命令
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "下载并更新GeoNames数据",
	Long:  `从GeoNames下载最新的地理位置数据并更新到数据库中。`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := downloadAndSaveData(); err != nil {
			logger.Logger.Error("下载数据失败", zap.Error(err))
			return
		}
		logger.Logger.Info("数据下载并保存成功")
	},
}

func downloadAndSaveData() error {
	// 下载数据
	locations, err := utils.DownloadGeoData()
	if err != nil {
		return fmt.Errorf("下载数据失败: %w", err)
	}

	// 获取存储实例
	storage := db.GetStorage()

	if err := storage.SaveLocations(locations); err != nil {
		return fmt.Errorf("批量保存数据失败: %w", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}

func main() {
	// 加载配置文件
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("加载配置文件失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	if err := logger.InitLogger(cfg.Log.Level, cfg.Log.Path); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化数据库连接
	if db := db.GetDB(); db == nil {
		logger.Logger.Error("初始化数据库连接失败")
		os.Exit(1)
	}
	defer func() {
		if err := db.CloseDB(); err != nil {
			logger.Logger.Error("关闭数据库连接失败", zap.Error(err))
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		logger.Logger.Error("执行命令失败", zap.Error(err))
		os.Exit(1)
	}
}
