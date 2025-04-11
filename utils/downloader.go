package utils

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/unxai/geonames-service/config"
	"github.com/unxai/geonames-service/logger"
	"github.com/unxai/geonames-service/models"
	"go.uber.org/zap"
)

const workerCount = 10 // 并发worker数量

// parseLocation 解析单行数据为Location结构
func parseLocation(line string) (models.Location, error) {
	fields := strings.Split(line, "\t")
	if len(fields) < 19 {
		return models.Location{}, fmt.Errorf("invalid field count")
	}

	geonameID, _ := strconv.Atoi(fields[0])
	lat, _ := strconv.ParseFloat(fields[4], 64)
	lon, _ := strconv.ParseFloat(fields[5], 64)
	pop, _ := strconv.Atoi(fields[14])
	elev, _ := strconv.Atoi(fields[15])

	return models.Location{
		GeonameID:        geonameID,
		Name:             fields[1],
		ASCII_Name:       fields[2],
		AlternateNames:   fields[3],
		Latitude:         lat,
		Longitude:        lon,
		FeatureClass:     fields[6],
		FeatureCode:      fields[7],
		CountryCode:      fields[8],
		Admin1Code:       fields[10],
		Admin2Code:       fields[11],
		Population:       pop,
		Elevation:        elev,
		TimeZone:         fields[17],
		ModificationDate: fields[18],
	}, nil
}

// DownloadGeoData 修改为返回数据而不是保存文件
func DownloadGeoData() ([]models.Location, error) {
	// 获取配置
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	// 检查本地缓存
	cacheFile := "data/allCountries.zip"
	if _, err = os.Stat(cacheFile); err == nil {
		// 如果缓存文件存在，直接使用缓存文件
		logger.Logger.Info("使用本地缓存文件")
		var data []byte
		data, err = os.ReadFile(cacheFile)
		if err != nil {
			return nil, fmt.Errorf("读取缓存文件失败: %w", err)
		}
		return parseZipData(data)
	}

	// 创建缓存目录
	if err = os.MkdirAll("data", 0755); err != nil {
		return nil, fmt.Errorf("创建缓存目录失败: %w", err)
	}

	logger.Logger.Info("开始下载数据文件")
	// 发起 HTTP 请求获取数据
	resp, err := http.Get(cfg.Download.URL)
	if err != nil {
		return nil, fmt.Errorf("下载数据文件失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应内容失败: %w", err)
	}

	// 保存到缓存文件
	if err := os.WriteFile(cacheFile, body, 0644); err != nil {
		logger.Logger.Warn("保存缓存文件失败", zap.Error(err))
	}

	return parseZipData(body)

}

// parseZipData 解析zip数据
func parseZipData(data []byte) ([]models.Location, error) {
	// 从内存中读取 zip 文件
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("解析zip文件失败: %w", err)
	}

	var locations []models.Location
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 使用更大的缓冲区来避免阻塞
	tasks := make(chan string, 10000)

	// 创建计数器用于跟踪进度
	var processedCount int32

	// 启动worker
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range tasks {
				location, err := parseLocation(line)
				if err != nil {
					logger.Logger.Warn("解析数据行失败", zap.Error(err))
					continue
				}

				mu.Lock()
				locations = append(locations, location)
				processedCount++
				if processedCount%1000 == 0 {
					logger.Logger.Info("处理进度", zap.Int32("processed", processedCount))
				}
				mu.Unlock()
			}
		}()
	}

	// 解析 zip 文件中的数据
	for _, file := range zipReader.File {
		if file.Name == "allCountries.txt" {
			rc, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("打开zip文件失败: %w", err)
			}

			// 使用scanner一次性处理数据
			scanner := bufio.NewScanner(rc)
			for scanner.Scan() {
				tasks <- scanner.Text()
			}

			// 检查scanner是否有错误
			if err := scanner.Err(); err != nil {
				rc.Close()
				return nil, fmt.Errorf("读取文件失败: %w", err)
			}

			rc.Close()
			break
		}
	}

	// 关闭任务channel
	close(tasks)

	// 等待所有worker完成
	wg.Wait()

	logger.Logger.Info("数据解析完成",
		zap.Int("total_locations", len(locations)))

	return locations, nil
}
