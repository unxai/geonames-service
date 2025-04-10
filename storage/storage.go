package storage

import (
	"github.com/unxai/geonames-service/models"
)

// Storage 定义了存储接口
type Storage interface {
	// SaveLocations 批量保存位置数据
	SaveLocations(locations []models.Location) error
}
