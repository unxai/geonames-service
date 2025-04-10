package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/unxai/geonames-service/logger"
	"github.com/unxai/geonames-service/models"
	"go.uber.org/zap"
)

// PostgresStorage 实现了 Storage 接口的 PostgreSQL 存储
type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage 创建一个新的 PostgreSQL 存储实例
func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

const batchSize = 1000 // 每批处理的数据量

// SaveLocations 批量保存位置数据
func (s *PostgresStorage) SaveLocations(locations []models.Location) error {
	if len(locations) == 0 {
		return nil
	}

	// 分批处理数据
	for i := 0; i < len(locations); i += batchSize {
		end := i + batchSize
		if end > len(locations) {
			end = len(locations)
		}
		batch := locations[i:end]

		// 开启事务
		tx, err := s.db.Begin()
		if err != nil {
			logger.Logger.Error("开启事务失败", zap.Error(err))
			return fmt.Errorf("开启事务失败: %w", err)
		}
		defer tx.Rollback()

		// 构建批量插入的值占位符
		valueStrings := make([]string, 0, len(batch))
		valueArgs := make([]interface{}, 0, len(batch)*9)
		for j, loc := range batch {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				j*9+1, j*9+2, j*9+3, j*9+4, j*9+5, j*9+6, j*9+7, j*9+8, j*9+9))
			valueArgs = append(valueArgs,
				loc.GeonameID, loc.Name, loc.ASCII_Name, loc.Latitude, loc.Longitude,
				loc.CountryCode, loc.Population, loc.FeatureClass, loc.FeatureCode)
		}

		// 构建完整的SQL语句
		sql := fmt.Sprintf(`
		INSERT INTO locations (
			geoname_id, name, ascii_name, latitude, longitude,
			country_code, population, feature_class, feature_code
		) VALUES %s
		ON CONFLICT (geoname_id) DO UPDATE SET
			name = EXCLUDED.name,
			ascii_name = EXCLUDED.ascii_name,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			country_code = EXCLUDED.country_code,
			population = EXCLUDED.population,
			feature_class = EXCLUDED.feature_class,
			feature_code = EXCLUDED.feature_code
		`, strings.Join(valueStrings, ","))

		// 执行批量插入
		_, err = tx.Exec(sql, valueArgs...)
		if err != nil {
			logger.Logger.Error("执行批量插入失败",
				zap.Int("batch_start", i),
				zap.Int("batch_size", len(batch)),
				zap.Error(err))
			return fmt.Errorf("执行批量插入失败: %w", err)
		}

		// 提交事务
		if err := tx.Commit(); err != nil {
			logger.Logger.Error("提交事务失败",
				zap.Int("batch_start", i),
				zap.Int("batch_size", len(batch)),
				zap.Error(err))
			return fmt.Errorf("提交事务失败: %w", err)
		}

		logger.Logger.Info("成功处理一批数据",
			zap.Int("batch_start", i),
			zap.Int("batch_size", len(batch)))
	}

	return nil
}
